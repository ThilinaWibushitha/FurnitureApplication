using System.Net.Http.Json;
using System.Text.Json;
using Microsoft.AspNetCore.Components;

namespace Furniture.Admin.Services
{
    public class AuthenticationService
    {
        private readonly HttpClient _httpClient;
        private readonly NavigationManager _navigationManager;

        public event Action<bool>? OnAuthenticationChanged;

        public bool IsAuthenticated { get; private set; }
        public string? UserName { get; private set; }
        public string? UserRole { get; private set; } // "admin" or "client"

        public AuthenticationService(HttpClient httpClient, NavigationManager navigationManager)
        {
            _httpClient = httpClient;
            _navigationManager = navigationManager;
        }

        public async Task<bool> LoginAsync(string username, string password, string role)
        {
            try
            {
                // For demo purposes, accept any non-empty credentials
                // In a real app, this would call your API
                if (!string.IsNullOrEmpty(username) && !string.IsNullOrEmpty(password))
                {
                    // Simulate API call
                    await Task.Delay(1000);

                    IsAuthenticated = true;
                    UserName = username;
                    UserRole = role;

                    OnAuthenticationChanged?.Invoke(true);

                    // Store in session storage (in real app, use secure tokens)
                    await StoreAuthDataAsync(username, role);

                    return true;
                }

                return false;
            }
            catch (Exception ex)
            {
                Console.WriteLine($"Login error: {ex.Message}");
                return false;
            }
        }

        public async Task LogoutAsync()
        {
            IsAuthenticated = false;
            UserName = null;
            UserRole = null;

            OnAuthenticationChanged?.Invoke(false);

            // Clear stored auth data
            await ClearAuthDataAsync();

            _navigationManager.NavigateTo("/login");
        }

        public async Task<bool> CheckAuthStatusAsync()
        {
            try
            {
                // Check if we have stored auth data
                var authData = await GetStoredAuthDataAsync();
                
                if (authData != null)
                {
                    UserName = authData.Username;
                    UserRole = authData.Role;
                    IsAuthenticated = true;
                    return true;
                }

                return false;
            }
            catch
            {
                return false;
            }
        }

        private async Task StoreAuthDataAsync(string username, string role)
        {
            var authData = new AuthData { Username = username, Role = role };
            var json = JsonSerializer.Serialize(authData);
            
            // In a real Blazor Server app, you might use session storage or a secure cookie
            // For demo, we'll use a simple approach
            await Task.CompletedTask;
        }

        private async Task<AuthData?> GetStoredAuthDataAsync()
        {
            // In a real app, retrieve from secure storage
            await Task.CompletedTask;
            return null; // For demo, always return null
        }

        private async Task ClearAuthDataAsync()
        {
            // Clear stored auth data
            await Task.CompletedTask;
        }

        private class AuthData
        {
            public string Username { get; set; } = string.Empty;
            public string Role { get; set; } = string.Empty;
        }
    }
}
