namespace StarFederation.Datastar.DependencyInjection

open System
open System.Text.Json
open Microsoft.AspNetCore.Http
open Microsoft.Extensions.DependencyInjection
open StarFederation.Datastar

module DatastarServices =

    let datastarServiceWithoutSignals (serviceCollection:IServiceCollection) =
        serviceCollection.AddHttpContextAccessor() |> ignore

        serviceCollection.AddScoped<IServerSentEventGenerator>(fun (svcPvd:IServiceProvider) ->
            let httpContext = svcPvd.GetService<IHttpContextAccessor>()
            ServerSentEventGenerator(httpContext)
            ) |> ignore

        serviceCollection

    let datastarService<'T when 'T:null and 'T :> IDatastarSignalsStore> (signalStoreDeserializer:string -> 'T) (serviceCollection:IServiceCollection) =
        serviceCollection.AddHttpContextAccessor() |> ignore

        serviceCollection.AddScoped<IServerSentEventGenerator>(fun (svcPvd:IServiceProvider) ->
            let httpContext = svcPvd.GetService<IHttpContextAccessor>()
            ServerSentEventGenerator(httpContext)
            ) |> ignore

        serviceCollection.AddScoped<IDatastarSignalsStore>(fun (svcPvd:IServiceProvider) ->
            let httpContextAccessor = svcPvd.GetService<IHttpContextAccessor>()
            let rawSignals = ServerSentEventHttpHandler.ReadRawSignalStore(httpContextAccessor.HttpContext.Request).GetAwaiter().GetResult()
            match rawSignals with
            | ValueSome rawSignals' -> signalStoreDeserializer(rawSignals')
            | ValueNone -> null
            ) |> ignore

        serviceCollection

[<System.Runtime.CompilerServices.Extension>]
type ServiceCollectionExtensionMethods() =

    [<System.Runtime.CompilerServices.Extension>]
    static member AddDatastarGenerator serviceCollection =
        DatastarServices.datastarServiceWithoutSignals serviceCollection

    [<System.Runtime.CompilerServices.Extension>]
    static member AddDatastarGenerator<'T when 'T:null and 'T :> IDatastarSignalsStore> serviceCollection =
        DatastarServices.datastarService<'T> JsonSerializer.Deserialize<'T> serviceCollection

    [<System.Runtime.CompilerServices.Extension>]
    static member AddDatastarGenerator<'T when 'T:null and 'T :> IDatastarSignalsStore> (serviceCollection:ServiceCollection, signalStoreDeserializer:Func<string, 'T>) =
        DatastarServices.datastarService<'T> signalStoreDeserializer.Invoke serviceCollection
