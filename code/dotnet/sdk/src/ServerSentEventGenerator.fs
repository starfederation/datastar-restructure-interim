namespace StarFederation.Datastar

open System.Text
open System.Text.Json
open System.Threading
open System.Threading.Tasks
open Microsoft.AspNetCore.Http
open Microsoft.Extensions.Primitives
open Microsoft.Net.Http.Headers

type ServerSentEventGenerator(httpContext:HttpContext) =
     let mutable inc = 0uL
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

     member private _.Context = httpContext

     new(httpContextAccessor:IHttpContextAccessor) =
        ServerSentEventGenerator(httpContextAccessor.HttpContext)

     member this.RenderFragment = ServerSideEvent.renderFragment this
     member this.RemoveFragment = ServerSideEvent.removeFragment this
     member this.PatchStore = ServerSideEvent.patchStore this
     member this.RemoveFromStore = ServerSideEvent.removeFromStore this
     member this.Redirect = ServerSideEvent.renderFragment this
     member this.Console = ServerSideEvent.console this

     static member ParseIncoming<'T> (httpRequest:HttpRequest) = task {
        let serializedDatastoreOptTask =
            match httpRequest.Method with
            | System.Net.WebRequestMethods.Http.Get ->
                Task.FromResult <|
                    match httpRequest.Query.TryGetValue("datastar") with
                    | true, json when json.Count > 0 -> ValueSome (json[0])
                    | _ -> ValueNone
            | System.Net.WebRequestMethods.Http.Post -> // ValueNone |> Task.FromResult
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

     with
     interface ISendEvent with
         member this.SendEvent(event:string) =
            let bytes = Encoding.UTF8.GetBytes(event)
            this.Context.Response.BodyWriter.WriteAsync(bytes).AsTask()

     interface IGetNextId with
        member this.GetNextId () = (Interlocked.Increment(&inc) - 1uL).ToString()