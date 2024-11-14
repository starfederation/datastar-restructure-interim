open System
open System.Text.Json
open System.Threading.Tasks
open Microsoft.AspNetCore.Builder
open Microsoft.AspNetCore.Http
open Microsoft.Extensions.Logging
open Microsoft.Extensions.Primitives
open Microsoft.Net.Http.Headers
open StarFederation.Datastar

type Store = { input: string; output: string; show: bool }

let storeDefaults = { input = ""; output = ""; show = true }
let indexPage = $"""<!doctype html>
<html>
<head>
    <title>F# + D* Example</title>
    <!--<script type="module" defer src="https://cdn.jsdelivr.net/gh/starfederation/datastar/bundles/datastar.js"></script>-->
    <script type="module" defer src="https://cdn.jsdelivr.net/npm/@sudodevnull/datastar"></script>
</head>
<body>
    <h2>F# + D* Example</h2>
    <main class="container" id="main" data-store='{storeDefaults |> JsonSerializer.Serialize}'>
        <input type="text" placeholder="Send to server..." data-model="input"/><br>
        <button data-on-click="$get('/get')">Send State Roundtrip</button><br>
        <button data-on-click="$get('/target')">Target HTML Element</button><br>
        <button data-on-click="$show=!$show">Toggle Feed</button><br>
        <div id="output" data-text="$output"></div>
        <div id="target"></div>
        <div data-show="$show">
            <span>Feed from server: </span>
            <span id="feed" data-on-load="$get('/feed')"></span>
        </div>

        <h5>Datastar Store</h5>
        <pre data-text="JSON.stringify(ctx.store(),null,2)"></pre>
    </main>
</body>
</html>"""

let addHeaders (headers: (string * string) list) (ctx:HttpContext) =
    let setHeader (name, content:string) =
        if ctx.Response.Headers.ContainsKey(name) |> not then
            ctx.Response.Headers.Add(name, StringValues(content))
    headers |> List.iter setHeader

module Handlers =

    let respondHtml (html:string) = RequestDelegate(fun ctx -> task {
        ctx.Response.Headers.Add(HeaderNames.ContentType, StringValues("text/html; charset=utf-8"))
        let bytes = System.Text.Encoding.UTF8.GetBytes(html)
        ctx.Response.ContentLength <- Nullable<int64>(bytes.LongLength)
        do! ctx.Response.BodyWriter.WriteAsync(bytes).AsTask() :> Task
        })

    let get = RequestDelegate(fun ctx -> task {
        let sse = ServerSentEventGenerator(ctx)
        let! store = ServerSentEventGenerator.ParseIncoming<Store>(ctx.Request)
        let store = store |> ValueOption.defaultValue storeDefaults
        let store = { store with output = $"Your input: {store.input}, is length {store.input.Length} long" }
        do! sse.RenderFragment { RenderFragmentOptions.defaults with MergeMode = UpsertAttributes } $@"<main id='main' data-store='{store |> JsonSerializer.Serialize}'></main>"
        })

    let target = RequestDelegate(fun ctx -> task {
        let sse = ServerSentEventGenerator(ctx)
        let today = DateTime.Now.ToString("%y-%M-%d %h:%m:%s")
        do! sse.RenderFragment RenderFragmentOptions.defaults $@"<div id='target'><b>{today}</b></div>"
        })

    let feed = RequestDelegate(fun ctx -> task {
        let sse = ServerSentEventGenerator(ctx)
        while not <| ctx.RequestAborted.IsCancellationRequested do
            let rand = Random.Shared.NextInt64()
            do! sse.RenderFragment RenderFragmentOptions.defaults $"<span id='feed'>{rand}</span>"
            do! Task.Delay (TimeSpan.FromSeconds 1)
        })

[<EntryPoint>]
let main args =
    let builder = WebApplication.CreateBuilder(args)

    let configureLogging (log:ILoggingBuilder) = log.ClearProviders().AddConsole()
    builder.Logging |> configureLogging |> ignore

    let app = builder.Build()

    app.MapGet("/", Handlers.respondHtml indexPage) |> ignore
    app.MapGet("/get", Handlers.get) |> ignore
    app.MapGet("/target", Handlers.target) |> ignore
    app.MapGet("/feed", Handlers.feed) |> ignore

    app.Run()

    0 // Exit code