namespace StarFederation.Datastar

open System
open System.Threading.Tasks

type IGetNextId = abstract member GetNextId: unit -> string
type ServerSideEventId =
    | Specified of string
    | Requested
module ServerSideEventId =
    let private getNextId (env:IGetNextId) = env.GetNextId()
    let value env id =
        match id with
        | Specified str -> str
        | Requested -> getNextId env

type FragmentMergeMode =
    | Morph | Inner | Outer | Prepend | Append
    | Before | After | UpsertAttributes
module internal FragmentMergeMode =
    let toString (this:FragmentMergeMode) =
        match this with
        | UpsertAttributes -> "upsertAttributes"
        | _ -> Utility.lowerFirstCharacter this

    let fromString (str:string) = Utility.unionCaseFromString<FragmentMergeMode> str [||]

    let defaultMergeMode =
        match fromString Defaults.DefaultFragmentMergeMode with
        | ValueSome def -> def
        | ValueNone -> failwith $"Default Merge Mode from {nameof Defaults}.fs, '{Defaults.DefaultFragmentMergeMode}' is unhandled. Add a union case to the FragmentMergeMode type"

type ConsoleMode =
    | Assert | Clear | Count | CountReset | Debug | Dir | Dirxml | Error
    | Group | GroupCollapsed | GroupEnd | Info | Log | Table | Time | TimeEnd
    | TimeLog | Trace | Warn
module internal ConsoleMode =
    let toString (this:ConsoleMode) = Utility.lowerFirstCharacter this

type DataStore = string
module DataStore =
    let value (dataStore:DataStore) = dataStore
    let tryCreate (dataStore:string) = Some (DataStore dataStore)
type DataStorePath = string
module DataStorePath =
    let value (dataStorePath:DataStorePath) = dataStorePath
    let tryCreate (dataStorePath:string) = Ok (DataStorePath dataStorePath)
type Selector = string
module Selector =
    let value (selector:Selector) = selector
    let tryCreate (selector:string) = Ok (Selector selector)
type IGetDataStore = abstract member GetDataStore: unit -> DataStore

type RenderFragmentOptions =
    { Selector: Selector voption
      MergeMode: FragmentMergeMode
      SettleDuration: TimeSpan
      UseViewTransition: bool }
module RenderFragmentOptions =
    let defaults = { Selector = ValueNone; MergeMode = FragmentMergeMode.defaultMergeMode; SettleDuration = Defaults.DefaultSettleTime; UseViewTransition = false }
type RemoveFragmentOptions =
    { SettleDuration: TimeSpan
      UseViewTransition: bool }
module RemoveFragmentOptions =
    let defaults = { SettleDuration = Defaults.DefaultSettleTime; UseViewTransition = false }

type ServerSideEventType = | Fragment | Signal | Remove | Redirect | Console
module ServerSideEventType =
    let toString this =
        match this with
        | Fragment -> "datastar-fragment"
        | Signal -> "datastar-signal"
        | Remove -> "datastar-remove"
        | Redirect -> "datastar-redirect"
        | Console -> "datastar-console"

type ISendEvent = abstract member SendEvent: string -> Task

type ServerSideEvent =
    { EventType: ServerSideEventType
      Id: ServerSideEventId
      Retry: TimeSpan voption
      DataLines: string[] }
module ServerSideEvent =
    let serializeEvent env (sseEvent:ServerSideEvent) =
        seq {
            $"event: {sseEvent.EventType |> ServerSideEventType.toString}"
            $"id: {sseEvent.Id |> ServerSideEventId.value env}"
            if (sseEvent.Retry |> ValueOption.isSome)
            then $"retry: {sseEvent.Retry |> ValueOption.get |> (_.Milliseconds)}"
            yield! sseEvent.DataLines |> Array.map (fun dataLine -> $"data: {dataLine}")
            ""; ""
        } |> String.concat "\n"

    let send env (sseEvent:ServerSideEvent) =
        // replace Id with
        let serializedEvent = sseEvent |> serializeEvent env
        let sendEvent (env:ISendEvent) (event:string) = env.SendEvent(event)
        sendEvent env serializedEvent

    let renderFragment env (options:RenderFragmentOptions) fragment =
        { EventType = Fragment
          Id = Requested
          Retry = ValueNone
          DataLines = seq {
                options.Selector          |> ValueOption.map (fun selector ->          $"selector {selector |> Selector.value}")
                options.MergeMode         |> (fun mergeMode ->               ValueSome $"merge {mergeMode |> FragmentMergeMode.toString}")
                options.SettleDuration    |> (fun settleDuration ->          ValueSome $"settle {settleDuration.Milliseconds}")
                options.UseViewTransition |> (fun useViewTransition ->       ValueSome $"vt {useViewTransition}")
                $"fragment %s{fragment}"  |> ValueSome
                } |> Seq.filter ValueOption.isSome |> Seq.map ValueOption.get |> Seq.toArray }
        |> send env

    let removeFragment env (options:RemoveFragmentOptions) selector =
        { EventType = Remove
          Id = Requested
          Retry = ValueNone
          DataLines = seq {
                $"selector {selector |> Selector.value}" |> ValueSome
                options.SettleDuration    |> (fun settleDuration ->          ValueSome $"settle {settleDuration.Milliseconds}")
                options.UseViewTransition |> (fun useViewTransition ->       ValueSome $"vt {useViewTransition}")
                } |> Seq.filter ValueOption.isSome |> Seq.map ValueOption.get |> Seq.toArray }
        |> send env

    let patchStore env onlyIfMissing data =
        { EventType = Signal
          Id = Requested
          Retry = ValueNone
          DataLines = [|
                if onlyIfMissing then "onlyIfMissing true"
                $"store {data |> DataStore.value}"
                |] }
        |> send env

    let removeFromStore env paths =
        let paths' = paths |> Seq.map DataStorePath.value |> String.concat " "
        { EventType = Remove
          Id = Requested
          Retry = ValueNone
          DataLines = [| $"selector {paths'}" |] }
        |> send env

    let redirect env url =
        { EventType = Redirect
          Id = Requested
          Retry = ValueNone
          DataLines = [| $"url %s{url}" |] }
        |> send env

    let console env mode message =
        { EventType = Console
          Id = Requested
          Retry = ValueNone
          DataLines = [| $"{mode |> ConsoleMode.toString} %s{message}" |] }
        |> send env
