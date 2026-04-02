using mobile.Services;
using System.Net.Http;
using System.Net.Sockets;
using System.Security.Authentication;

namespace mobile.Pages;

public partial class SettingsPage : ContentPage
{
	private readonly FrontendSettings frontendSettings;

	public SettingsPage(FrontendSettings frontendSettings)
	{
		InitializeComponent();
		this.frontendSettings = frontendSettings;
	}

	protected override void OnAppearing()
	{
		base.OnAppearing();
		FrontendUrlEntry.Text = frontendSettings.FrontendUrl;
		UpdateSavedValue(frontendSettings.HasFrontendUrl ? frontendSettings.FrontendUrl : "Not configured");
	}

	private async void OnSaveClicked(object? sender, EventArgs e)
	{
		if (!frontendSettings.TrySave(FrontendUrlEntry.Text, out var normalizedUrl, out var errorMessage))
		{
			StatusLabel.TextColor = Color.FromArgb("#B42318");
			StatusLabel.Text = errorMessage;
			await DisplayAlertAsync("Invalid URL", errorMessage, "OK");
			return;
		}

		FrontendUrlEntry.Text = normalizedUrl;
		UpdateSavedValue(normalizedUrl);
		StatusLabel.TextColor = Color.FromArgb("#1B6E3C");
		StatusLabel.Text = "Saved. Switch back to the UpSnap tab to load the site.";
	}

	private async void OnTestClicked(object? sender, EventArgs e)
	{
		if (!frontendSettings.TryGetNormalizedUri(FrontendUrlEntry.Text, out var uri, out var normalizedUrl, out var errorMessage))
		{
			StatusLabel.TextColor = Color.FromArgb("#B42318");
			StatusLabel.Text = errorMessage;
			await DisplayAlertAsync("Invalid URL", errorMessage, "OK");
			return;
		}

		StatusLabel.TextColor = Color.FromArgb("#666666");
		StatusLabel.Text = "Testing connectivity from this device...";

		var (success, message) = await TestUrlAsync(uri!);
		StatusLabel.TextColor = Color.FromArgb(success ? "#1B6E3C" : "#B42318");
		StatusLabel.Text = message;

		if (success)
		{
			FrontendUrlEntry.Text = normalizedUrl;
		}
	}

	private void OnClearClicked(object? sender, EventArgs e)
	{
		frontendSettings.Clear();
		FrontendUrlEntry.Text = string.Empty;
		UpdateSavedValue("Not configured");
		StatusLabel.TextColor = Color.FromArgb("#666666");
		StatusLabel.Text = "Saved frontend URL cleared.";
	}

	private void UpdateSavedValue(string value)
	{
		SavedUrlLabel.Text = value;
	}

	private async void OnDoneClicked(object? sender, EventArgs e)
	{
		await Navigation.PopModalAsync();
	}

	private static async Task<(bool Success, string Message)> TestUrlAsync(Uri uri)
	{
		try
		{
			return await Task.Run(async () =>
			{
				using var httpClient = new HttpClient
				{
					Timeout = TimeSpan.FromSeconds(12),
				};

				using var request = new HttpRequestMessage(HttpMethod.Get, uri);
				using var response = await httpClient.SendAsync(request, HttpCompletionOption.ResponseHeadersRead).ConfigureAwait(false);

				if (response.IsSuccessStatusCode)
				{
					return (true, $"Connection test passed. The device reached {uri} and got HTTP {(int)response.StatusCode}.");
				}

				return (false, $"The device reached {uri}, but the server returned HTTP {(int)response.StatusCode} ({response.ReasonPhrase}).");
			}).ConfigureAwait(false);
		}
		catch (HttpRequestException ex) when (ex.InnerException is SocketException socketEx)
		{
			return socketEx.SocketErrorCode switch
			{
				SocketError.HostNotFound or SocketError.NoData =>
					(false, $"DNS lookup failed for {uri.Host}. Android emulators often do not resolve local names such as my-pi.lan the same way Windows does. Try the server's LAN IP or a public DNS name."),
				SocketError.NetworkUnreachable or SocketError.HostUnreachable =>
					(false, $"The emulator has no route to {uri.Host}. If this address is on a VPN or overlay network, the Android emulator usually cannot use the Windows-only route. Try a normal LAN IP that the emulator can reach."),
				SocketError.ConnectionRefused =>
					(false, $"The emulator reached {uri.Host}:{uri.Port}, but the connection was refused. Check that the UpSnap frontend is listening on that interface and port, and that the firewall allows it."),
				SocketError.TimedOut =>
					(false, $"Connection to {uri} timed out. The emulator could not reach the server in time. This usually means routing, firewall, or reverse-proxy reachability issues."),
				_ =>
					(false, $"Connection test failed for {uri} with network error {socketEx.SocketErrorCode}. If this works on Windows only, treat it as an emulator DNS or routing issue first.")
			};
		}
		catch (HttpRequestException ex) when (ex.InnerException is AuthenticationException)
		{
			return (false, $"TLS failed for {uri}. If this is an https URL, confirm the certificate is trusted by Android and matches the hostname.");
		}
		catch (TaskCanceledException)
		{
			return (false, $"Connection to {uri} timed out. The device could not reach the server in time.");
		}
		catch (Exception ex)
		{
			if (string.Equals(ex.GetType().FullName, "Android.OS.NetworkOnMainThreadException", StringComparison.Ordinal))
			{
				return (false, $"Android blocked the connectivity test on the UI thread. Rebuild with the latest code and try again. If the next error mentions DNS or routing, the emulator still cannot reach {uri.Host} from its network.");
			}

			return (false, $"Connection test failed for {uri}: {ex.Message}");
		}
	}
}