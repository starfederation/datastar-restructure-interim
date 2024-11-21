namespace StarFederation.Datastar.DependencyInjection

open System
open System.IO
open System.Text
open System.Text.Json
open System.Text.Json.Serialization
open System.Threading.Tasks
open Microsoft.AspNetCore.Http
open Microsoft.Extensions.Primitives
open Microsoft.Net.Http.Headers
open StarFederation.Datastar
open StarFederation.Datastar.Utility

type IServerSentEventHandler =
    interface
        inherit ISendServerEvent
        inherit IReadRawSignalsStore
    end

type ServerSentEventGenerator (handler:IServerSentEventHandler) =
    new (httpContext:HttpContext) =
        ServerSentEventGenerator (ServerSentEventHttpHandler httpContext)
    new (httpContextAccessor:IHttpContextAccessor) =
        ServerSentEventGenerator httpContextAccessor.HttpContext

    member private _.Handler = handler

    member this.ReadSignalsAsync<'T when 'T:null and 'T :> IDatastarSignalsStore> (deserializeSignalStore:Func<string, 'T>) : Task<'T> = task {
        let! storeVOpt = this.Handler.ReadRawSignalStore()
        match storeVOpt with
        | ValueNone -> return! Task.FromException<'T>(exn "Failed to parse Signal Store")
        | ValueSome serializedSignalStore -> return deserializeSignalStore.Invoke(serializedSignalStore)
        }
    member this.ReadSignalsAsync<'T when 'T:null and 'T :> IDatastarSignalsStore> () : Task<'T> =
        this.ReadSignalsAsync (fun str -> JsonSerializer.Deserialize<'T>(str))

    with
    interface IServerSentEventGenerator with
        member this.MergeFragments(fragment) = ServerSentEvent.mergeFragments this.Handler fragment
        member this.MergeFragments(fragment, options) = ServerSentEvent.mergeFragmentsWithOptions options.AsOptions this.Handler fragment
        member this.MergeSignals(dataSignals) = ServerSentEvent.mergeSignals this.Handler Consts.DefaultMergeSignalsOnlyIfMissing dataSignals
        member this.MergeSignals(dataSignals, onlyIfMissing): Task = ServerSentEvent.mergeSignals this.Handler onlyIfMissing dataSignals
        member this.MergeSignals(dataSignals, options:ServerSentEventOptions): Task = ServerSentEvent.mergeSignalsWithOptions options.AsOptions this.Handler Consts.DefaultMergeSignalsOnlyIfMissing dataSignals
        member this.MergeSignals(dataSignals, onlyIfMissing, options:ServerSentEventOptions): Task = ServerSentEvent.mergeSignalsWithOptions options.AsOptions this.Handler onlyIfMissing dataSignals
        member this.RemoveFragments(selector) = ServerSentEvent.removeFragments this.Handler selector
        member this.RemoveFragments(selector: Selector, options: ServerSentEventRemoveFragmentsOptions) = ServerSentEvent.removeFragmentsWithOptions options.AsOptions this.Handler selector
        member this.RemoveSignals(paths) = ServerSentEvent.removeSignals this.Handler paths
        member this.RemoveSignals(paths, options) = ServerSentEvent.removeSignalsWithOptions options.AsOptions this.Handler paths
        member this.ExecuteScript(script) = ServerSentEvent.executeScript this.Handler script
        member this.ExecuteScript(script, options) = ServerSentEvent.executeScriptWithOptions options.AsOptions this.Handler script

and ServerSentEventHttpHandler(httpContext:HttpContext) =
    do
        let setHeader (ctx:HttpContext) (name, content:string) =
            if ctx.Response.Headers.ContainsKey(name) |> not then
                ctx.Response.Headers.Add(name, StringValues(content))
        [
           ("Cache-Control", "no-cache")
           ("Connection", "keep-alive")
           (HeaderNames.ContentType, "text/event-stream")
        ] |> Seq.iter (setHeader httpContext)
        httpContext.Response.StartAsync().GetAwaiter().GetResult()
        httpContext.Response.Body.FlushAsync().GetAwaiter().GetResult()

    member _.HttpContext = httpContext

    static member ReadRawSignalStore (httpRequest:HttpRequest) : ValueTask<string voption> =
        match httpRequest.Method with
        | System.Net.WebRequestMethods.Http.Get ->
            match httpRequest.Query.TryGetValue(Consts.DatastarKey) with
            | true, json when json.Count > 0 -> ValueSome (json[0])
            | _ -> ValueNone
            |> ValueTask.FromResult
        | System.Net.WebRequestMethods.Http.Post ->
            task {
                try
                    use readResult = new StreamReader(httpRequest.BodyReader.AsStream())
                    let! str = readResult.ReadToEndAsync()
                    return (ValueSome str)
                with _ -> return ValueNone
                } |> ValueTask<string voption>
        | _ -> ValueNone |> ValueTask.FromResult

    with
    interface IServerSentEventHandler

    interface ISendServerEvent with
        member this.SendServerEvent(event:string) =
            let bytes = Encoding.UTF8.GetBytes(event)
            this.HttpContext.Response.BodyWriter.WriteAsync(bytes).AsTask()

    interface IReadRawSignalsStore with
        member this.ReadRawSignalStore () = ServerSentEventHttpHandler.ReadRawSignalStore(this.HttpContext.Request)

/// Server Sent Generator Interface
and IServerSentEventGenerator =
    abstract MergeFragments: fragment:string -> Task
    abstract MergeFragments: fragment:string * options:ServerSentEventMergeFragmentsOptions -> Task
    abstract RemoveFragments: selector:Selector -> Task
    abstract RemoveFragments: selector:Selector * options:ServerSentEventRemoveFragmentsOptions -> Task
    abstract MergeSignals: dataSignals:IDatastarSignalsStore -> Task
    abstract MergeSignals: dataSignals:IDatastarSignalsStore * options:ServerSentEventOptions -> Task
    abstract MergeSignals: dataSignals:IDatastarSignalsStore * onlyIfMissing:bool -> Task
    abstract MergeSignals: dataSignals:IDatastarSignalsStore * onlyIfMissing:bool * options:ServerSentEventOptions -> Task
    abstract RemoveSignals: paths:DataSignalPath seq -> Task
    abstract RemoveSignals: paths:DataSignalPath seq * options:ServerSentEventOptions -> Task
    abstract ExecuteScript: script:string -> Task
    abstract ExecuteScript: script:string * options:ServerSentEventExecuteScriptOptions -> Task

and ServerSentEventMergeFragmentsOptions() =
    member val EventId:string = null with get, set
    member val Retry = Consts.DefaultSseRetryDuration with get, set
    member val MergeMode = MergeFragmentsOptions.defaults.MergeMode with get, set
    member val Selector:string = null with get, set
    member val SettleDuration = MergeFragmentsOptions.defaults.SettleDuration with get, set
    member val UseViewTransition = MergeFragmentsOptions.defaults.UseViewTransition with get, set
    member this.AsOptions : MergeFragmentsOptions =
        { Selector = this.Selector |> ValueOption.fromNullable
          MergeMode = this.MergeMode
          SettleDuration = this.SettleDuration
          UseViewTransition = this.UseViewTransition
          EventId = this.EventId |> ValueOption.fromNullable
          Retry = this.Retry }

and ServerSentEventRemoveFragmentsOptions() =
    member val EventId:string = null with get, set
    member val Retry = Consts.DefaultSseRetryDuration with get, set
    member val SettleDuration = MergeFragmentsOptions.defaults.SettleDuration with get, set
    member val UseViewTransition = MergeFragmentsOptions.defaults.UseViewTransition with get, set
    member this.AsOptions : RemoveFragmentsOptions =
        { SettleDuration = this.SettleDuration
          UseViewTransition = this.UseViewTransition
          EventId = this.EventId |> ValueOption.fromNullable
          Retry = this.Retry }

and ServerSentEventOptions() =
    member val EventId:string = null with get, set
    member val Retry = Consts.DefaultSseRetryDuration with get, set
    member this.AsOptions : EventOptions =
        { EventId = this.EventId |> ValueOption.fromNullable
          Retry = this.Retry }

and ServerSentEventExecuteScriptOptions() =
    member val AutoRemove:bool = Consts.DefaultExecuteScriptAutoRemove with get, set
    member val Attributes:string[] = [| Consts.DefaultExecuteScriptAttributes |] with get, set
    member val EventId:string = null with get, set
    member val Retry = Consts.DefaultSseRetryDuration with get, set
    member this.AsOptions : ExecuteScriptOptions =
        { AutoRemove = this.AutoRemove
          Attributes = this.Attributes
          EventId = this.EventId |> ValueOption.fromNullable
          Retry = this.Retry }
