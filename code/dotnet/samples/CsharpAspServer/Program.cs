using System.Text.Json;
using System.Text.Json.Serialization;
using Microsoft.Extensions.FileProviders;
using StarFederation.Datastar;
using StarFederation.Datastar.DependencyInjection;

namespace CsharpAspServer;

public record DataSignalsStore : IDatastarSignalsStore
{
    [JsonPropertyName("input")]
    [JsonIgnore(Condition = JsonIgnoreCondition.WhenWritingNull)]
    public string? Input { get; init; } = null;

    [JsonPropertyName("output")]
    [JsonIgnore(Condition = JsonIgnoreCondition.WhenWritingNull)]
    public string? Output { get; init; } = null;

    [JsonPropertyName("show")]
    [JsonIgnore(Condition = JsonIgnoreCondition.WhenWritingNull)]
    public bool? Show { get; init; } = null;

    public string Serialize() => JsonSerializer.Serialize(this);
}

public static class Program
{
    public static void Main(string[] args)
    {
        WebApplicationBuilder builder = WebApplication.CreateBuilder(args);
        builder.Services.AddDatastarGenerator<DataSignalsStore>();

        WebApplication app = builder.Build();
        app.UseDefaultFiles(new DefaultFilesOptions
        {
            FileProvider = new PhysicalFileProvider(Path.Combine(builder.Environment.ContentRootPath, "..", "Shared", "wwwroot")),
        });
        app.UseStaticFiles(new StaticFileOptions
        {
            FileProvider = new PhysicalFileProvider(Path.Combine(builder.Environment.ContentRootPath, "..", "Shared", "wwwroot")),
        });

        app.MapGet("/language/{lang:required}", (string lang, IServerSentEventGenerator sse) => sse.MergeFragments($"""<span id="language">{lang}</span>"""));
        app.MapGet("/patch", async (IServerSentEventGenerator sse, IDatastarSignalsStore dsStore) =>
        {
            DataSignalsStore signalsStore = (dsStore as DataSignalsStore) ?? throw new InvalidCastException("Unknown Datastore passed");
            DataSignalsStore mergeSignalsStore = new() { Output = $"Patched Output: {signalsStore.Input}" };
            await sse.MergeSignals(mergeSignalsStore);
        });
        app.MapGet("/target", async (IServerSentEventGenerator sse) =>
        {
            string today = DateTime.Now.ToString("%y-%M-%d %h:%m:%s");
            await sse.MergeFragments($"""<div id='target'><span id='date'><b>{today}</b><button data-on-click="$get('/removeDate')">Remove</button></span></div>""");
        });
        app.MapGet("/removeDate", (IServerSentEventGenerator sse) => sse.RemoveFragments("#date"));
        app.MapGet("/feed", async (IHttpContextAccessor acc, IServerSentEventGenerator sse, CancellationToken ct) =>
        {
            try
            {
                while (!ct.IsCancellationRequested)
                {
                    long rand = Random.Shared.NextInt64(1000000000000000000, 5999999999999999999);
                    await sse.MergeFragments($"<span id='feed'>{rand}</span>");
                    await Task.Delay(TimeSpan.FromSeconds(1), ct);
                }
            }
            finally
            {
                try { acc?.HttpContext?.Connection.RequestClose(); } catch { }
            }
        });

        app.Run();
    }
}
