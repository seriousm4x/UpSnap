namespace mobile.Pages;

public sealed class MainTabsPage : TabbedPage
{
	public MainTabsPage(BrowserPage browserPage, SettingsPage settingsPage)
	{
		Title = "UpSnap";
		Children.Add(browserPage);
		Children.Add(settingsPage);
	}
}