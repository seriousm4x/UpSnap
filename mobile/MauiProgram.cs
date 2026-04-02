using Microsoft.Extensions.Logging;

namespace mobile;

public static class MauiProgram
{
	public static MauiApp CreateMauiApp()
	{
		var builder = MauiApp.CreateBuilder();
		builder
			.UseMauiApp<App>()
			.ConfigureFonts(fonts =>
			{
				fonts.AddFont("OpenSans-Regular.ttf", "OpenSansRegular");
				fonts.AddFont("OpenSans-Semibold.ttf", "OpenSansSemibold");
			});

#if ANDROID
		builder.ConfigureMauiHandlers(handlers =>
		{
			Microsoft.Maui.Handlers.WebViewHandler.Mapper.AppendToMapping("UpSnapAndroidWebView", (handler, _) =>
			{
				var settings = handler.PlatformView.Settings;
				settings.JavaScriptEnabled = true;
				settings.DomStorageEnabled = true;
				settings.MixedContentMode = Android.Webkit.MixedContentHandling.CompatibilityMode;
			});
		});
#endif

		builder.Services.AddSingleton<Services.FrontendSettings>();
		builder.Services.AddSingleton<Pages.BrowserPage>();
		builder.Services.AddTransient<Pages.SettingsPage>();

#if DEBUG
		builder.Logging.AddDebug();
#endif

		return builder.Build();
	}
}
