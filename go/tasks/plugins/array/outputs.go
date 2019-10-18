package array

import (
	"context"
	"fmt"

	arrayCore "github.com/lyft/flyteplugins/go/tasks/plugins/array/core"

	"github.com/lyft/flyteplugins/go/tasks/pluginmachinery/workqueue"
	"github.com/lyft/flytestdlib/bitarray"
	"github.com/lyft/flytestdlib/promutils"
	"github.com/lyft/flytestdlib/storage"

	"github.com/lyft/flyteplugins/go/tasks/pluginmachinery/io"
	"github.com/lyft/flyteplugins/go/tasks/plugins/array/errorcollector"

	"github.com/lyft/flyteplugins/go/tasks/pluginmachinery/ioutils"

	"github.com/lyft/flyteidl/gen/pb-go/flyteidl/core"
	pluginCore "github.com/lyft/flyteplugins/go/tasks/pluginmachinery/core"
)

var (
	nilLiteral = &core.Literal{
		Value: &core.Literal_Scalar{
			Scalar: &core.Scalar{
				Value: &core.Scalar_NoneType{
					NoneType: &core.Void{},
				},
			},
		},
	}
)

// Represents an indexed work queue that aggregates outputs of sub tasks.
type OutputAssembler struct {
	workqueue.IndexedWorkQueue
}

func (o OutputAssembler) Queue(id workqueue.WorkItemID, item *outputAssembleItem) error {
	return o.IndexedWorkQueue.Queue(id, item)
}

type outputAssembleItem struct {
	outputPaths io.OutputFilePaths
	varNames    []string
	finalPhases bitarray.CompactArray
	dataStore   *storage.DataStore
}

type assembleOutputsWorker struct {
}

func (w assembleOutputsWorker) Process(ctx context.Context, workItem workqueue.WorkItem) (workqueue.WorkStatus, error) {
	i := workItem.(*outputAssembleItem)

	outputReaders, err := ConstructOutputReaders(ctx, i.dataStore, i.outputPaths.GetOutputPrefixPath(), int(i.finalPhases.ItemsCount))
	if err != nil {
		return workqueue.WorkStatusNotDone, err
	}

	finalOutputs := &core.LiteralMap{
		Literals: map[string]*core.Literal{},
	}
	for idx, subTaskPhaseIdx := range i.finalPhases.GetItems() {
		existingPhase := pluginCore.Phases[subTaskPhaseIdx]
		if existingPhase.IsSuccess() {
			output, executionError, err := outputReaders[idx].Read(ctx)
			if err != nil {
				return workqueue.WorkStatusNotDone, err
			}

			if executionError == nil && output != nil {
				appendSubTaskOutput(finalOutputs, output, int64(i.finalPhases.ItemsCount))
				continue
			}
		}

		// TODO: Do we need the names of the outputs in the literalMap here?
		appendEmptyOutputs(finalOutputs, i.varNames)
	}

	ow := ioutils.NewRemoteFileOutputWriter(ctx, i.dataStore, i.outputPaths)
	if err = ow.Put(ctx, ioutils.NewInMemoryOutputReader(finalOutputs, nil)); err != nil {
		return workqueue.WorkStatusNotDone, err
	}

	return workqueue.WorkStatusSucceeded, nil
}

func appendOneItem(outputs *core.LiteralMap, varName string, literal *core.Literal, expectedSize int64) {
	existingVal, found := outputs.Literals[varName]
	var list *core.LiteralCollection
	if found {
		list = existingVal.GetCollection()
	} else {
		list = &core.LiteralCollection{
			Literals: make([]*core.Literal, 0, expectedSize),
		}

		existingVal = &core.Literal{
			Value: &core.Literal_Collection{
				Collection: list,
			},
		}
	}

	list.Literals = append(list.Literals, literal)
	outputs.Literals[varName] = existingVal
}

func appendSubTaskOutput(outputs *core.LiteralMap, subTaskOutput *core.LiteralMap,
	expectedSize int64) {

	for key, val := range subTaskOutput.GetLiterals() {
		appendOneItem(outputs, key, val, expectedSize)
	}
}

func appendEmptyOutputs(outputs *core.LiteralMap, vars []string) {
	for _, varName := range vars {
		appendOneItem(outputs, varName, nilLiteral, 1)
	}
}

// Assembles a single outputs.pb that contain all the outputs of the subtasks and write them to the final OutputWriter.
// This step can potentially be expensive (hence the metrics) and why it's offloaded to a background process.
func AssembleFinalOutputs(ctx context.Context, assemblyQueue OutputAssembler, tCtx pluginCore.TaskExecutionContext,
	terminalPhase arrayCore.Phase, state *arrayCore.State) (*arrayCore.State, error) {

	// Otherwise, run the data catalog steps - create and submit work items to the catalog processor,
	// build input readers
	workItemID := tCtx.TaskExecutionMetadata().GetTaskExecutionID().GetGeneratedName()
	w, found, err := assemblyQueue.Get(workItemID)
	if err != nil {
		return nil, err
	}

	// If the work item is not found in the queue, it's either never queued, was evicted, or we are recovering from a
	// reboot. Add it again.
	if !found {
		taskTemplate, err := tCtx.TaskReader().Read(ctx)
		if err != nil {
			return nil, err
		}

		outputVariables := taskTemplate.GetInterface().GetOutputs()
		if outputVariables == nil || outputVariables.GetVariables() == nil {
			// If the task has no outputs, bail early.
			state = state.SetPhase(terminalPhase, 0)
			return state, nil
		}

		varNames := make([]string, 0, len(outputVariables.GetVariables()))

		err = assemblyQueue.Queue(workItemID, &outputAssembleItem{
			varNames:    varNames,
			finalPhases: state.GetArrayStatus().Detailed,
			outputPaths: tCtx.OutputWriter(),
			dataStore:   tCtx.DataStore(),
		})

		if err != nil {
			return nil, err
		}

		w, found, err = assemblyQueue.Get(workItemID)
		if err != nil {
			return nil, err
		}

		if !found {
			return nil, fmt.Errorf("couldn't find work item [%v] after immediately adding it", workItemID)
		}
	}

	switch w.Status() {
	case workqueue.WorkStatusSucceeded:
		or := ioutils.NewRemoteFileOutputReader(ctx, tCtx.DataStore(), tCtx.OutputWriter(), tCtx.MaxDatasetSizeBytes())
		if err = tCtx.OutputWriter().Put(ctx, or); err != nil {
			return nil, err
		}

		state = state.SetPhase(terminalPhase, 0)
	case workqueue.WorkStatusFailed:
		state = state.SetExecutionErr(&core.ExecutionError{
			Message: w.Error().Error(),
		})

		state = state.SetPhase(arrayCore.PhaseRetryableFailure, 0)
	}

	return state, nil
}

type assembleErrorsWorker struct {
	maxErrorMessageLength int
}

func (a assembleErrorsWorker) Process(ctx context.Context, workItem workqueue.WorkItem) (workqueue.WorkStatus, error) {
	w := workItem.(*outputAssembleItem)
	outputReaders, err := ConstructOutputReaders(ctx, w.dataStore, w.outputPaths.GetOutputPrefixPath(), int(w.finalPhases.ItemsCount))
	if err != nil {
		return workqueue.WorkStatusNotDone, err
	}

	ec := errorcollector.NewErrorMessageCollector()
	for idx, subTaskPhaseIdx := range w.finalPhases.GetItems() {
		existingPhase := pluginCore.Phases[subTaskPhaseIdx]
		if existingPhase.IsFailure() {
			isError, err := outputReaders[idx].IsError(ctx)

			if err != nil {
				return workqueue.WorkStatusNotDone, err
			}

			if isError {
				executionError, err := outputReaders[idx].ReadError(ctx)
				if err != nil {
					return workqueue.WorkStatusNotDone, err
				}

				ec.Collect(idx, executionError.String())
			}
		}
	}

	msg := ""
	if ec.Length() > 0 {
		msg = ec.Summary(a.maxErrorMessageLength)
	}

	ow := ioutils.NewRemoteFileOutputWriter(ctx, w.dataStore, w.outputPaths)
	if err = ow.Put(ctx, ioutils.NewInMemoryOutputReader(nil, &io.ExecutionError{
		ExecutionError: &core.ExecutionError{
			Code:     "",
			Message:  msg,
			ErrorUri: "",
		},
		IsRecoverable: false,
	})); err != nil {
		return workqueue.WorkStatusNotDone, err
	}

	return workqueue.WorkStatusSucceeded, nil
}

func NewOutputAssembler(workQueueConfig workqueue.Config, scope promutils.Scope) (OutputAssembler, error) {
	q, err := workqueue.NewIndexedWorkQueue("output", assembleOutputsWorker{}, workQueueConfig, scope)
	if err != nil {
		return OutputAssembler{}, err
	}

	return OutputAssembler{
		IndexedWorkQueue: q,
	}, nil
}

func NewErrorAssembler(maxErrorMessageLength int, workQueueConfig workqueue.Config, scope promutils.Scope) (OutputAssembler, error) {
	q, err := workqueue.NewIndexedWorkQueue("error", assembleErrorsWorker{
		maxErrorMessageLength: maxErrorMessageLength,
	}, workQueueConfig, scope)

	if err != nil {
		return OutputAssembler{}, err
	}

	return OutputAssembler{
		IndexedWorkQueue: q,
	}, nil
}
