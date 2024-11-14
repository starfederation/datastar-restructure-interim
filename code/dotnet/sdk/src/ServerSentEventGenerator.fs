namespace StarFederation.Datastar

open System.Text
open System.Text.Json
open System.Threading.Tasks
open Microsoft.AspNetCore.Http
open Microsoft.Extensions.Primitives
open Microsoft.Net.Http.Headers

type IServerSentEventContext =
    interface
        inherit ISendEvent
    end

type ServerSentEventGenerator (context:IServerSentEventContext) =
    new (httpContext:HttpContext) =
        ServerSentEventGenerator (ServerSentEventHttpContext httpContext)
    new (httpContextAccessor:IHttpContextAccessor) =
        ServerSentEventGenerator httpContextAccessor.HttpContext

    member private _.Context = context

    member this.RenderFragment = ServerSentEvent.renderFragment this.Context EventOptions.defaults
    member this.RemoveFragment = ServerSentEvent.removeFragment this.Context EventOptions.defaults
    member this.PatchStore = ServerSentEvent.patchStore this.Context EventOptions.defaults
    member this.RemoveFromStore = ServerSentEvent.removeFromStore this.Context EventOptions.defaults
    member this.Redirect = ServerSentEvent.redirect this.Context EventOptions.defaults
    member this.Console = ServerSentEvent.console this.Context EventOptions.defaults

    static member ParseIncoming (httpRequest:HttpRequest) = task {
        let serializedDatastoreOptTask =
            match httpRequest.Method with
            | System.Net.WebRequestMethods.Http.Get ->
                Task.FromResult <|
                    match httpRequest.Query.TryGetValue(Consts.datastarQueryKey) with
                    | true, json when json.Count > 0 -> ValueSome (json[0])
                    | _ -> ValueNone
            | System.Net.WebRequestMethods.Http.Post ->
                task {
                    try
                        let! readResult = httpRequest.BodyReader.ReadAsync().AsTask()
                        let bytes = readResult.Buffer
                        let str = Encoding.UTF8.GetString(&bytes)
                        return (ValueSome str)
                    with _ -> return ValueNone
                    }
            | _ -> ValueNone |> Task.FromResult
        let! serializedDatastoreOpt = serializedDatastoreOptTask
        return
            try serializedDatastoreOpt |> ValueOption.map (fun serializedDatastore -> JsonSerializer.Deserialize<'T>(serializedDatastore))
            with _ -> ValueNone
            }

    static member ParseIncomingAsync<'T when 'T:null and 'T :> IDatastarStore> (httpRequest:HttpRequest) : Task<'T> = task {
        let! storeVOpt = ServerSentEventGenerator.ParseIncoming(httpRequest)
        return
            match storeVOpt with
            | ValueSome store -> store
            | ValueNone -> null
     }

    with
    interface IServerSentEventGenerator with
        member this.RenderFragment(fragment) = ServerSentEvent.renderFragment this.Context EventOptions.defaults RenderFragmentOptions.defaults fragment
        member this.RenderFragment(fragment, options) = ServerSentEvent.renderFragment this.Context EventOptions.defaults options.AsOptions fragment
        member this.RenderFragment(fragment, options, eventOptions) = ServerSentEvent.renderFragment this.Context eventOptions.AsOptions options.AsOptions fragment
        member this.RemoveFragment(selector) = ServerSentEvent.removeFragment this.Context EventOptions.defaults RemoveFragmentOptions.defaults selector
        member this.RemoveFragment(selector, options) = ServerSentEvent.removeFragment this.Context EventOptions.defaults options.AsOptions selector
        member this.RemoveFragment(selector, options, eventOptions) = ServerSentEvent.removeFragment this.Context eventOptions.AsOptions { SettleDuration = options.SettleDuration; UseViewTransition = options.UseViewTransition } selector
        member this.PatchStore(dataStorePatch) = ServerSentEvent.patchStore this.Context EventOptions.defaults false dataStorePatch
        member this.PatchStore(dataStorePatch:string, eventOptions:ServerSentEventOpts) = ServerSentEvent.patchStore this.Context eventOptions.AsOptions false dataStorePatch
        member this.PatchStore(dataStorePatch:string, onlyIfMissing:bool) = ServerSentEvent.patchStore this.Context EventOptions.defaults onlyIfMissing dataStorePatch
        member this.PatchStore(dataStorePatch:string, onlyIfMissing:bool, eventOptions:ServerSentEventOpts) = ServerSentEvent.patchStore this.Context eventOptions.AsOptions onlyIfMissing dataStorePatch
        member this.RemoveFromStore(dataStorePaths) = ServerSentEvent.removeFromStore this.Context EventOptions.defaults dataStorePaths
        member this.RemoveFromStore(dataStorePaths, eventOptions) = ServerSentEvent.removeFromStore this.Context eventOptions.AsOptions dataStorePaths
        member this.Redirect(uri) = ServerSentEvent.redirect this.Context EventOptions.defaults uri
        member this.Redirect(uri, eventOptions) = ServerSentEvent.redirect this.Context eventOptions.AsOptions uri
        member this.Console(consoleMode, message) = ServerSentEvent.console this.Context EventOptions.defaults consoleMode message
        member this.Console(consoleMode, eventOptions, message) = ServerSentEvent.console this.Context eventOptions.AsOptions consoleMode message

and ServerSentEventHttpContext(httpContext:HttpContext) =
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

    member _.Context = httpContext

    with
    interface IServerSentEventContext

    interface ISendEvent with
        member this.SendEvent(event:string) =
            let bytes = Encoding.UTF8.GetBytes(event)
            this.Context.Response.BodyWriter.WriteAsync(bytes).AsTask()

and IServerSentEventGenerator =
    abstract member RenderFragment: fragment:string -> Task
    abstract member RenderFragment: fragment:string * options:RenderFragmentOpts -> Task
    abstract member RenderFragment: fragment:string * options:RenderFragmentOpts * eventOptions:ServerSentEventOpts -> Task
    abstract member RemoveFragment: selector:Selector -> Task
    abstract member RemoveFragment: selector:Selector * options:RemoveFragmentOpts -> Task
    abstract member RemoveFragment: selector:Selector * options:RemoveFragmentOpts * eventOptions:ServerSentEventOpts -> Task
    abstract member PatchStore: dataStore:string -> Task
    abstract member PatchStore: dataStore:string * eventOptions:ServerSentEventOpts -> Task
    abstract member PatchStore: dataStore:string * onlyIfMissing:bool -> Task
    abstract member PatchStore: dataStore:string * onlyIfMissing:bool * eventOptions:ServerSentEventOpts -> Task
    abstract member RemoveFromStore: paths:DataStorePath seq -> Task
    abstract member RemoveFromStore: paths:DataStorePath seq * eventOptions:ServerSentEventOpts -> Task
    abstract member Redirect: url:string -> Task
    abstract member Redirect: url:string * eventOptions:ServerSentEventOpts -> Task
    abstract member Console: consoleMode:ConsoleMode * message:string -> Task
    abstract member Console: consoleMode:ConsoleMode * eventOptions:ServerSentEventOpts * message:string -> Task

and RenderFragmentOpts() =
    member val MergeMode = RenderFragmentOptions.defaults.MergeMode with get, set
    member val Selector = null with get, set
    member val SettleDuration = RenderFragmentOptions.defaults.SettleDuration with get, set
    member val UseViewTransition = RenderFragmentOptions.defaults.UseViewTransition with get, set
    member this.AsOptions =
        { MergeMode = this.MergeMode
          Selector = Utility.ValueOption.fromNullable this.Selector
          SettleDuration = this.SettleDuration
          UseViewTransition = this.UseViewTransition }

and RemoveFragmentOpts() =
    member val SettleDuration = RenderFragmentOptions.defaults.SettleDuration with get, set
    member val UseViewTransition = RenderFragmentOptions.defaults.UseViewTransition with get, set
    member this.AsOptions =
        { SettleDuration = this.SettleDuration
          UseViewTransition = this.UseViewTransition }

and ServerSentEventOpts() =
    member val EventId = null with get, set
    member val Retry = Consts.DefaultSseSendRetry with get, set
    member this.AsOptions =
        { EventId = Utility.ValueOption.fromNullable this.EventId
          Retry = this.Retry }
