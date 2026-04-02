using System.ComponentModel;
using System.Runtime.CompilerServices;
using System.Text.Json;
using Microsoft.Maui.Storage;

namespace mobile.Services;

public sealed class FrontendSettings : INotifyPropertyChanged
{
	private const string FrontendUrlPreferenceKey = "frontend_url";
	private readonly string settingsFilePath;
	private string frontendUrl;

	public FrontendSettings()
	{
		settingsFilePath = Path.Combine(FileSystem.Current.AppDataDirectory, "frontend-settings.json");
		frontendUrl = LoadFrontendUrl();
	}

	public event PropertyChangedEventHandler? PropertyChanged;

	public string FrontendUrl
	{
		get => frontendUrl;
		private set
		{
			if (frontendUrl == value)
			{
				return;
			}

			frontendUrl = value;
			OnPropertyChanged();
			OnPropertyChanged(nameof(HasFrontendUrl));
		}
	}

	public bool HasFrontendUrl => !string.IsNullOrWhiteSpace(FrontendUrl);

	public Uri? GetUriOrNull()
	{
		if (!Uri.TryCreate(FrontendUrl, UriKind.Absolute, out var uri))
		{
			return null;
		}

		return IsSupportedScheme(uri) ? uri : null;
	}

	public bool TrySave(string? candidate, out string normalizedUrl, out string errorMessage)
	{
		if (!TryNormalize(candidate, out normalizedUrl, out errorMessage))
		{
			return false;
		}

		try
		{
			PersistSettings(normalizedUrl);
			Preferences.Default.Remove(FrontendUrlPreferenceKey);
		}
		catch (Exception ex) when (ex is IOException or UnauthorizedAccessException)
		{
			errorMessage = "The frontend URL could not be saved on this device.";
			return false;
		}

		FrontendUrl = normalizedUrl;
		return true;
	}

	public void Clear()
	{
		if (File.Exists(settingsFilePath))
		{
			File.Delete(settingsFilePath);
		}

		Preferences.Default.Remove(FrontendUrlPreferenceKey);
		FrontendUrl = string.Empty;
	}

	private string LoadFrontendUrl()
	{
		if (TryLoadFromFile(out var storedUrl))
		{
			return storedUrl;
		}

		var legacyValue = NormalizeStoredValue(Preferences.Default.Get(FrontendUrlPreferenceKey, string.Empty));
		if (string.IsNullOrEmpty(legacyValue))
		{
			return string.Empty;
		}

		PersistSettings(legacyValue);
		Preferences.Default.Remove(FrontendUrlPreferenceKey);
		return legacyValue;
	}

	private void PersistSettings(string normalizedUrl)
	{
		Directory.CreateDirectory(FileSystem.Current.AppDataDirectory);
		var document = new FrontendSettingsDocument
		{
			FrontendUrl = normalizedUrl,
		};

		var json = JsonSerializer.Serialize(document);
		File.WriteAllText(settingsFilePath, json);
	}

	private bool TryLoadFromFile(out string storedUrl)
	{
		storedUrl = string.Empty;

		if (!File.Exists(settingsFilePath))
		{
			return false;
		}

		try
		{
			var json = File.ReadAllText(settingsFilePath);
			var document = JsonSerializer.Deserialize<FrontendSettingsDocument>(json);
			storedUrl = NormalizeStoredValue(document?.FrontendUrl);
			return !string.IsNullOrEmpty(storedUrl);
		}
		catch (Exception ex) when (ex is IOException or JsonException or UnauthorizedAccessException)
		{
			storedUrl = string.Empty;
			return false;
		}
	}

	private static bool IsSupportedScheme(Uri uri)
	{
		return string.Equals(uri.Scheme, Uri.UriSchemeHttp, StringComparison.OrdinalIgnoreCase)
			|| string.Equals(uri.Scheme, Uri.UriSchemeHttps, StringComparison.OrdinalIgnoreCase);
	}

	private static string NormalizeStoredValue(string? value)
	{
		return TryNormalize(value, out var normalizedUrl, out _) ? normalizedUrl : string.Empty;
	}

	private static bool TryNormalize(string? candidate, out string normalizedUrl, out string errorMessage)
	{
		normalizedUrl = string.Empty;
		errorMessage = string.Empty;

		if (string.IsNullOrWhiteSpace(candidate))
		{
			errorMessage = "Enter the full frontend URL.";
			return false;
		}

		candidate = candidate.Trim();

		if (!Uri.TryCreate(candidate, UriKind.Absolute, out var uri) || !IsSupportedScheme(uri))
		{
			errorMessage = "Use a full http or https URL, for example http://192.168.1.50:5173.";
			return false;
		}

		normalizedUrl = uri.AbsoluteUri.TrimEnd('/');
		return true;
	}

	private void OnPropertyChanged([CallerMemberName] string? propertyName = null)
	{
		PropertyChanged?.Invoke(this, new PropertyChangedEventArgs(propertyName));
	}

	private sealed class FrontendSettingsDocument
	{
		public string? FrontendUrl { get; set; }
	}
}