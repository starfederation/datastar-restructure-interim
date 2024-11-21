namespace StarFederation.Datastar

open System
open System.Collections.Generic
open System.Text.RegularExpressions

[<AllowNullLiteral>]
type IDatastarSignalsStore =
    abstract Serialize : unit -> string

type DataSignalPath = string
module DataSignalPath =
    let value = id
    let create (dataSignalPath:string) = DataSignalPath dataSignalPath
    let tryCreate (dataSignalPath:string) = ValueSome (create dataSignalPath)

type Selector = string
module Selector =
    let value = id
    let regex = Regex(@"[#.][-_]?[_a-zA-Z]+(?:\w|\\.)*|(?<=\s+|^)(?:\w+|\*)|\[[^\s""'=<>`]+?(?<![~|^$*])([~|^$*]?=(?:['""].*['""]|[^\s""'=<>`]+))?\]|:[\w-]+(?:\(.*\))?", RegexOptions.Compiled)
    let create (selectorString:string) = Selector selectorString
    let tryCreate (selector:string) =
        match regex.IsMatch selector with
        | true -> ValueSome (create selector)
        | false -> ValueNone

type MergeFragmentsOptions =
    { Selector: Selector voption
      MergeMode: FragmentMergeMode
      SettleDuration: TimeSpan
      UseViewTransition: bool
      EventId: string voption
      Retry: TimeSpan }
module MergeFragmentsOptions =
    let defaults =
        { Selector = ValueNone
          MergeMode = Consts.DefaultFragmentMergeMode
          SettleDuration = Consts.DefaultSettleDuration
          UseViewTransition = Consts.DefaultFragmentsUseViewTransitions
          EventId = ValueNone
          Retry = Consts.DefaultSseRetryDuration }

type RemoveFragmentsOptions =
    { SettleDuration: TimeSpan
      UseViewTransition: bool
      EventId: string voption
      Retry: TimeSpan }
module RemoveFragmentsOptions =
    let defaults =
        { SettleDuration = Consts.DefaultSettleDuration
          UseViewTransition = Consts.DefaultFragmentsUseViewTransitions
          EventId = ValueNone
          Retry = Consts.DefaultSseRetryDuration }

type ExecuteScriptOptions =
    { AutoRemove: bool
      Attributes: string[]
      EventId: string voption
      Retry: TimeSpan }
module ExecuteScriptOptions =
    let defaults =
        { AutoRemove = Consts.DefaultExecuteScriptAutoRemove
          Attributes = [| Consts.DefaultExecuteScriptAttributes |]
          EventId = ValueNone
          Retry = Consts.DefaultSseRetryDuration }


type EventOptions = { EventId: string voption; Retry: TimeSpan }
module EventOptions =
    let defaults = { EventId = ValueNone; Retry = Consts.DefaultSseRetryDuration }

