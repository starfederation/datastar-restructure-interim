namespace StarFederation.Datastar

open System
open System.Text.RegularExpressions
open System.Threading.Tasks

type IDatastarStore = interface end
type DataStorePath = string
module DataStorePath =
    let value (dataStorePath:DataStorePath) = dataStorePath
    let create (dataStorePath:string) = DataStorePath dataStorePath
    let tryCreate (dataStorePath:string) = ValueSome (create dataStorePath)

type Selector = string
module Selector =
    let value (selector:Selector) = selector
    let regex = Regex(@"[#.][-_]?[_a-zA-Z]+(?:\w|\\.)*|(?<=\s+|^)(?:\w+|\*)|\[[^\s""'=<>`]+?(?<![~|^$*])([~|^$*]?=(?:['""].*['""]|[^\s""'=<>`]+))?\]|:[\w-]+(?:\(.*\))?", RegexOptions.Compiled)
    let create (selectorString:string) = Selector selectorString
    let tryCreate (selector:string) =
        match regex.IsMatch selector with
        | true -> ValueSome (create selector)
        | false -> ValueNone

type RenderFragmentOptions =
    { Selector: Selector voption
      MergeMode: FragmentMergeMode
      SettleDuration: TimeSpan
      UseViewTransition: bool }
module RenderFragmentOptions =
    let defaults = { Selector = ValueNone; MergeMode = Default; SettleDuration = Consts.DefaultSettleTime; UseViewTransition = Consts.DefaultUseViewTransition }

type RemoveFragmentOptions = { SettleDuration: TimeSpan; UseViewTransition: bool }
module RemoveFragmentOptions =
    let defaults = { SettleDuration = Consts.DefaultSettleTime; UseViewTransition = Consts.DefaultUseViewTransition }

type EventOptions = { EventId: string voption; Retry: TimeSpan voption }
module EventOptions =
    let defaults = { EventId = ValueNone; Retry = ValueNone }

type ISendEvent = abstract member SendEvent: string -> Task

type ServerSentEvent =
    { EventType: ServerSentEventType
      Id: string voption
      Retry: TimeSpan voption
      DataLines: string[] }
module ServerSentEvent =
    let serializeEvent sseEvent =
        seq {
            $"event: {sseEvent.EventType |> Consts.ServerSentEventType.toString}"

            if sseEvent.Id |> ValueOption.isSome
            then $"id: {sseEvent.Id |> ValueOption.get}"

            if (sseEvent.Retry |> ValueOption.isSome)
            then $"retry: {sseEvent.Retry |> ValueOption.get |> (_.Milliseconds)}"

            yield! sseEvent.DataLines |> Array.map (fun dataLine -> $"data: {dataLine}")

            ""; ""
        } |> String.concat "\n"

    let send env sseEvent =
        let serializedEvent = sseEvent |> serializeEvent
        let sendEvent (env:ISendEvent) (event:string) = env.SendEvent(event)
        sendEvent env serializedEvent

    let renderFragment env eventOptions options fragment =
        { EventType = Fragment
          Id = eventOptions.EventId
          Retry = eventOptions.Retry
          DataLines = [|
            if options.Selector |> ValueOption.isSome then $"{Consts.dataSelector} {options.Selector |> ValueOption.get |> Selector.value}"
            $"{Consts.dataMerge} {options.MergeMode |> Consts.FragmentMergeMode.toString}"
            $"{Consts.dataSettleDuration} {options.SettleDuration.Milliseconds}"
            if options.UseViewTransition then $"{Consts.dataUseViewTransition} true"
            $"{Consts.dataFragment} %s{fragment}"
            |] }
        |> send env

    let removeFragment env eventOptions options selector =
        { EventType = Remove
          Id = eventOptions.EventId
          Retry = eventOptions.Retry
          DataLines = [|
            $"{Consts.dataSelector} {selector |> Selector.value}"
            $"{Consts.dataSettleDuration} {options.SettleDuration.Milliseconds}"
            if options.UseViewTransition then $"{Consts.dataUseViewTransition} true"
            |] }
        |> send env

    let patchStore env eventOptions onlyIfMissing data =
        { EventType = Signal
          Id = eventOptions.EventId
          Retry = eventOptions.Retry
          DataLines = [|
            if onlyIfMissing then $"{Consts.dataOnlyIfMissing} true"
            $"{Consts.dataStore} %s{data}"
            |] }
        |> send env

    let removeFromStore env eventOptions paths =
        let paths' = paths |> Seq.map DataStorePath.value |> String.concat " "
        { EventType = Remove
          Id = eventOptions.EventId
          Retry = eventOptions.Retry
          DataLines = [| $"{Consts.dataSelector} {paths'}" |] }
        |> send env

    let redirect env eventOptions url =
        { EventType = Redirect
          Id = eventOptions.EventId
          Retry = eventOptions.Retry
          DataLines = [| $"{Consts.dataUrl} %s{url}" |] }
        |> send env

    let console env eventOptions mode message =
        { EventType = ServerSentEventType.Console
          Id = eventOptions.EventId
          Retry = eventOptions.Retry
          DataLines = [| $"{mode |> Consts.ConsoleMode.toString} %s{message}" |] }
        |> send env
