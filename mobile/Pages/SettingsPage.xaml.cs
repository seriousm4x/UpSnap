using mobile.Services;

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
}