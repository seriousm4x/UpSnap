namespace mobile;

public partial class App : Application
{
	private readonly Pages.BrowserPage browserPage;

	public App(Pages.BrowserPage browserPage)
	{
		InitializeComponent();
		this.browserPage = browserPage;
	}

	protected override Window CreateWindow(IActivationState? activationState)
	{
		return new Window(new NavigationPage(browserPage));
	}
}