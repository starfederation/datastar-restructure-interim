namespace StarFederation.Datastar

open System
open System.Threading.Tasks

type ISendServerEvent = abstract SendServerEvent: string -> Task
type IReadRawSignalsStore = abstract ReadRawSignalStore: unit -> ValueTask<string voption>

type ServerSentEvent =
    { EventType: EventType
      Id: string voption
      Retry: TimeSpan
      DataLines: string[] }
module ServerSentEvent =
    let serializeEvent sseEvent =
        seq {
            $"event: {sseEvent.EventType |> Consts.EventType.toString}"

            if sseEvent.Id |> ValueOption.isSome
            then $"id: {sseEvent.Id |> ValueOption.get}"

            if (sseEvent.Retry <> Consts.DefaultSseRetryDuration)
            then $"retry: {sseEvent.Retry.Milliseconds}"

            yield! sseEvent.DataLines |> Array.map (fun dataLine -> $"data: {dataLine}")

            ""; ""
        } |> String.concat "\n"


    let send env sseEvent =
        let serializedEvent = sseEvent |> serializeEvent
        let sendEvent (env:ISendServerEvent) (event:string) = env.SendServerEvent(event)
        sendEvent env serializedEvent

    let mergeFragmentsWithOptions (options:MergeFragmentsOptions) env (fragment:string) =
        { EventType = MergeFragments
          Id = options.EventId
          Retry = options.Retry
          DataLines = [|
            if (options.Selector |> ValueOption.isSome) then $"{Consts.DatastarDatalineSelector} {options.Selector |> ValueOption.get |> Selector.value}"
            if (options.MergeMode <> Consts.DefaultFragmentMergeMode) then $"{Consts.DatastarDatalineMergeMode} {options.MergeMode |> Consts.FragmentMergeMode.toString}"
            if (options.SettleDuration <> Consts.DefaultSettleDuration) then $"{Consts.DatastarDatalineSettleDuration} {options.SettleDuration.Milliseconds}"
            if (options.UseViewTransition <> Consts.DefaultFragmentsUseViewTransitions) then $"{Consts.DatastarDatalineUseViewTransition} {options.UseViewTransition |> Utility.toLower}"
            yield! (fragment |> Utility.splitLine |> Seq.map (fun fragmentLine -> $"{Consts.DatastarDatalineFragments} %s{fragmentLine}"))
            |] }
        |> send env
    let mergeFragments env = mergeFragmentsWithOptions MergeFragmentsOptions.defaults env

    let removeFragmentsWithOptions (options:RemoveFragmentsOptions) env selector =
        { EventType = RemoveFragments
          Id = options.EventId
          Retry = options.Retry
          DataLines = [|
            $"{Consts.DatastarDatalineSelector} {selector |> Selector.value}"
            if (options.SettleDuration <> Consts.DefaultSettleDuration) then $"{Consts.DatastarDatalineSettleDuration} {options.SettleDuration.Milliseconds}"
            if (options.UseViewTransition <> Consts.DefaultFragmentsUseViewTransitions) then $"{Consts.DatastarDatalineUseViewTransition} {options.UseViewTransition |> Utility.toLower}"
            |] }
        |> send env
    let removeFragments env = removeFragmentsWithOptions RemoveFragmentsOptions.defaults env

    let mergeSignalsWithOptions options env onlyIfMissing (mergeSignalData:IDatastarSignalsStore) : Task =
        { EventType = MergeSignals
          Id = options.EventId
          Retry = options.Retry
          DataLines = [|
            if (onlyIfMissing <> Consts.DefaultMergeSignalsOnlyIfMissing) then $"{Consts.DatastarDatalineOnlyIfMissing} {onlyIfMissing |> Utility.toLower}"
            yield! (mergeSignalData.Serialize() |> Utility.splitLine |> Seq.map (fun dataLine -> $"{Consts.DatastarDatalineSignals} %s{dataLine}"))
            |] }
       |> send env
    let mergeSignals env = mergeSignalsWithOptions EventOptions.defaults env

    let removeSignalsWithOptions options env paths =
        let paths' = paths |> Seq.map DataSignalPath.value |> String.concat " "
        { EventType = RemoveSignals
          Id = options.EventId
          Retry = options.Retry
          DataLines = [| $"{Consts.DatastarDatalineSelector} {paths'}" |] }
        |> send env
    let removeSignals env = removeSignalsWithOptions EventOptions.defaults env

    let executeScriptWithOptions (options:ExecuteScriptOptions) env script =
        { EventType = ExecuteScript
          Id = options.EventId
          Retry = options.Retry
          DataLines = [|
            if (options.AutoRemove <> Consts.DefaultExecuteScriptAutoRemove) then $"{Consts.DatastarDatalineAutoRemove} {options.AutoRemove |> Utility.toLower}"
            if (not <| Seq.forall2 (=) options.Attributes [| Consts.DefaultExecuteScriptAttributes |] ) then $"{Consts.DefaultExecuteScriptAttributes} {options.AutoRemove |> Utility.toLower}"
            yield! script |> Utility.splitLine |> Seq.map (fun scriptLine -> $"{Consts.DatastarDatalineScript} %s{scriptLine}")
          |] }
        |> send env
    let executeScript env = executeScriptWithOptions ExecuteScriptOptions.defaults env
