using System.ComponentModel;
using Microsoft.Extensions.DependencyInjection;
using mobile.Services;

namespace mobile.Pages;

public partial class BrowserPage : ContentPage
{
	private readonly FrontendSettings frontendSettings;
	private readonly IServiceProvider serviceProvider;

	public BrowserPage(FrontendSettings frontendSettings, IServiceProvider serviceProvider)
	{
		InitializeComponent();
		this.frontendSettings = frontendSettings;
		this.serviceProvider = serviceProvider;
		this.frontendSettings.PropertyChanged += OnFrontendSettingsPropertyChanged;
	}

	protected override void OnAppearing()
	{
		base.OnAppearing();
		RefreshWebView();
	}

	private void OnFrontendSettingsPropertyChanged(object? sender, PropertyChangedEventArgs e)
	{
		if (e.PropertyName is nameof(FrontendSettings.FrontendUrl) or nameof(FrontendSettings.HasFrontendUrl))
		{
			MainThread.BeginInvokeOnMainThread(RefreshWebView);
		}
	}

	private void RefreshWebView()
	{
		var uri = frontendSettings.GetUriOrNull();
		EmptyStateLayout.IsVisible = uri is null;
		FrontendWebView.IsVisible = uri is not null;

		if (uri is null)
		{
			LoadingOverlay.IsVisible = false;
			FrontendWebView.Source = null;
			return;
		}

		var currentUrl = (FrontendWebView.Source as UrlWebViewSource)?.Url;
		if (!string.Equals(currentUrl, uri.AbsoluteUri, StringComparison.OrdinalIgnoreCase))
		{
			FrontendWebView.Source = uri.AbsoluteUri;
		}
	}

	private async void OnShowUrlClicked(object? sender, EventArgs e)
	{
		var message = frontendSettings.HasFrontendUrl
			? frontendSettings.FrontendUrl
			: "Set the frontend URL in Settings.";

		await DisplayAlertAsync("Frontend URL", message, "OK");
	}

	private async void OnOpenSettingsClicked(object? sender, EventArgs e)
	{
		var settingsPage = serviceProvider.GetRequiredService<SettingsPage>();
		await Navigation.PushModalAsync(new NavigationPage(settingsPage));
	}

	private void OnRefreshClicked(object? sender, EventArgs e)
	{
		if (FrontendWebView.IsVisible)
		{
			FrontendWebView.Reload();
		}
	}

	private void OnWebViewNavigating(object? sender, WebNavigatingEventArgs e)
	{
		LoadingOverlay.IsVisible = true;
	}

	private async void OnWebViewNavigated(object? sender, WebNavigatedEventArgs e)
	{
		LoadingOverlay.IsVisible = false;

		if (e.Result != WebNavigationResult.Success)
		{
			var message = $"The WebView could not load {frontendSettings.FrontendUrl}. Check the address in Settings and confirm the site is reachable from this device.";

#if ANDROID
			message += " If this is an http URL, Android also needs cleartext traffic enabled. If this is an https URL, check that the certificate is valid on the device.";
#endif

			await DisplayAlertAsync("Frontend unreachable", message, "OK");
		}
	}
}