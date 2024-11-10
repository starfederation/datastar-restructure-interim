open System
open System.Text.Json
open System.Threading
open System.Threading.Tasks
open Falco
open Falco.HttpHandler
open Falco.Extensions
open Falco.Routing
open Falco.HostBuilder
open Microsoft.AspNetCore.Http

type Store = { input: string; output: string; show: bool }
let store = { input = ""; output = ""; show = true }
let indexPage = $"""<!doctype html>
<html>
<head>
    <title>F# + Datastar Example</title>
    <script type="module" defer src="https://cdn.jsdelivr.net/npm/@sudodevnull/datastar"></script>
</head>
<body>
    <h2>Node/Express + Datastar Example</h2>
    <main class="container" id="main" data-store='{store |> JsonSerializer.Serialize}'>
        <input type="text" placeholder="Send to server..." data-model="input"/><br>
        <button data-on-click="$$get('/get')">Send State Roundtrip</button><br>
        <button data-on-click="$$get('/target')">Target HTML Element</button><br>
        <button data-on-click="$show=!$show">Toggle Feed</button><br>
        <div id="output" data-text="$output"></div>
        <div id="target"></div>
        <div data-show="$show">
            <span>Feed from server: </span>
            <span id="feed" data-on-load="$$get('/feed')"></span>
        </div>

        <h5>Datastar Store</h5>
        <pre data-text="JSON.stringify(ctx.store(),null,2)"></pre>
    </main>
</body>
</html>"""

module Response =
    let withDatastarHeaders =
        Response.withHeaders [ ("Cache-Control", "no-cache"); ("Connection", "keep-alive") ]
        >> Response.withContentType "text/event-stream"
    let ofHtmlFragment fragment =
        withDatastarHeaders
        >> Response.ofHtmlString (DatastarFragment.serialize fragment)
    let ofHtmlFragmentStream fragments = (fun ctx -> task {
        withDatastarHeaders ctx |> ignore
        for fragment in fragments do
            do! (ctx.Response.BodyWriter.WriteAsync (DatastarFragment.serializeToBytes fragment)).AsTask() :> Task
        })

module Handlers =
    let get = (fun ctx ->
        let form = Request.getQuery ctx
        let store = form.GetString("datastar", "") |> JsonSerializer.Deserialize
        let store' = { store with output = $"""Your input: {store.input}, is {String.length store.input} long""" }
        Response.ofHtmlFragment { Merge = true; Data = $@"<main id='main' data-store='{store' |> JsonSerializer.Serialize}'></main>" } ctx
        )

    let target = (fun ctx ->
        let today = DateTime.Now.ToString("%y-%M-%d %h:%m:%s")
        Response.ofHtmlFragment { Merge = false; Data = $@"<div id='target'><b>{today}</b></div>" } ctx
        )

    let feed = (fun (ctx:HttpContext) ->
        Response.ofHtmlFragmentStream (seq {
            while not <| ctx.RequestAborted.IsCancellationRequested do
                let rand = Random.Shared.NextInt64()
                yield { Merge = false; Data = $"<span id='feed'>{rand}</span>" }
                Thread.Sleep 1000
            }) ctx :> Task
        )

[<EntryPoint>]
let main args =
    webHost args {
        endpoints [
            get "/" (Response.ofHtmlString indexPage)
            get "/get" Handlers.get
            get "/target" Handlers.target
            get "/feed" Handlers.feed
            ]
    }
    0