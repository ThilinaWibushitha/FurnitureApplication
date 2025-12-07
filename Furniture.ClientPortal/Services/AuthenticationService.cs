using System.Net.Http.Json;
using System.Text.Json;
using Microsoft.AspNetCore.Components;
using Furniture.ClientPortal.Models;

namespace Furniture.ClientPortal.Services
{
    public class AuthenticationService : IAuthenticationService
    {
        private readonly HttpClient _httpClient;
        private readonly NavigationManager _navigationManager;

        public event Action<bool>? OnAuthenticationChanged;

        public bool IsAuthenticated { get; private set; }
        public string? UserName { get; private set; }
        public UserRole CurrentUserRole { get; private set; } = UserRole.Client;
        public User? CurrentUser { get; private set; }

        public AuthenticationService(HttpClient httpClient, NavigationManager navigationManager)
        {
            _httpClient = httpClient;
            _navigationManager = navigationManager;
        }

        // Main Admin Authentication
        public async Task<AuthResponse?> MainAdminLoginAsync(MainAdminLoginRequest request)
        {
            try
            {
                var response = await _httpClient.PostAsJsonAsync("api/auth/main-admin/login", request);
                if (response.IsSuccessStatusCode)
                {
                    var authResponse = await response.Content.ReadFromJsonAsync<AuthResponse>();
                    await SetAuthState(authResponse);
                    return authResponse;
                }
                return null;
            }
            catch (Exception ex)
            {
                Console.WriteLine($"Main admin login error: {ex.Message}");
                return null;
            }
        }

        public async Task<bool> MainAdminPasswordResetRequestAsync(string email)
        {
            try
            {
                var response = await _httpClient.PostAsJsonAsync("api/auth/main-admin/reset-request", new { email });
                return response.IsSuccessStatusCode;
            }
            catch (Exception ex)
            {
                Console.WriteLine($"Main admin password reset request error: {ex.Message}");
                return false;
            }
        }

        // Regular Admin Authentication
        public async Task<AuthResponse?> AdminLoginAsync(AdminLoginRequest request)
        {
            try
            {
                var response = await _httpClient.PostAsJsonAsync("api/auth/admin/login", request);
                if (response.IsSuccessStatusCode)
                {
                    var authResponse = await response.Content.ReadFromJsonAsync<AuthResponse>();
                    await SetAuthState(authResponse);
                    return authResponse;
                }
                return null;
            }
            catch (Exception ex)
            {
                Console.WriteLine($"Admin login error: {ex.Message}");
                return null;
            }
        }

        public async Task<bool> AdminPasswordResetRequestAsync(string username)
        {
            try
            {
                var response = await _httpClient.PostAsJsonAsync("api/auth/admin/reset-request", new { username });
                return response.IsSuccessStatusCode;
            }
            catch (Exception ex)
            {
                Console.WriteLine($"Admin password reset request error: {ex.Message}");
                return false;
            }
        }

        // Client Authentication
        public async Task<AuthResponse?> ClientLoginAsync(LoginRequest request)
        {
            try
            {
                var response = await _httpClient.PostAsJsonAsync("api/auth/login", request);
                if (response.IsSuccessStatusCode)
                {
                    var authResponse = await response.Content.ReadFromJsonAsync<AuthResponse>();
                    await SetAuthState(authResponse);
                    return authResponse;
                }
                return null;
            }
            catch (Exception ex)
            {
                Console.WriteLine($"Client login error: {ex.Message}");
                return null;
            }
        }

        public async Task<bool> ClientPasswordResetRequestAsync(string email)
        {
            try
            {
                var response = await _httpClient.PostAsJsonAsync("api/auth/client/reset-request", new { email });
                return response.IsSuccessStatusCode;
            }
            catch (Exception ex)
            {
                Console.WriteLine($"Client password reset request error: {ex.Message}");
                return false;
            }
        }

        public async Task<bool> ConfirmPasswordResetAsync(PasswordResetConfirm confirm)
        {
            try
            {
                var response = await _httpClient.PostAsJsonAsync("api/auth/reset-confirm", confirm);
                return response.IsSuccessStatusCode;
            }
            catch (Exception ex)
            {
                Console.WriteLine($"Password reset confirm error: {ex.Message}");
                return false;
            }
        }

        // User Management (Main Admin only)
        public async Task<User?> CreateUserAsync(CreateUserRequest request)
        {
            try
            {
                var response = await _httpClient.PostAsJsonAsync("api/users", request);
                if (response.IsSuccessStatusCode)
                {
                    return await response.Content.ReadFromJsonAsync<User>();
                }
                return null;
            }
            catch (Exception ex)
            {
                Console.WriteLine($"Create user error: {ex.Message}");
                return null;
            }
        }

        public async Task<User?> CreateAdminAsync(CreateAdminRequest request)
        {
            try
            {
                var response = await _httpClient.PostAsJsonAsync("api/users/create-admin", request);
                if (response.IsSuccessStatusCode)
                {
                    return await response.Content.ReadFromJsonAsync<User>();
                }
                return null;
            }
            catch (Exception ex)
            {
                Console.WriteLine($"Create admin error: {ex.Message}");
                return null;
            }
        }

        public async Task<bool> DeleteUserAsync(long userId)
        {
            try
            {
                var response = await _httpClient.DeleteAsync($"api/users/{userId}");
                return response.IsSuccessStatusCode;
            }
            catch (Exception ex)
            {
                Console.WriteLine($"Delete user error: {ex.Message}");
                return false;
            }
        }

        public async Task<List<User>> GetUsersAsync()
        {
            try
            {
                return await _httpClient.GetFromJsonAsync<List<User>>("api/users") ?? new List<User>();
            }
            catch (Exception ex)
            {
                Console.WriteLine($"Get users error: {ex.Message}");
                return new List<User>();
            }
        }

        // Password Change Management
        public async Task<bool> RequestPasswordChangeAsync(PasswordChangeRequest request)
        {
            try
            {
                var response = await _httpClient.PostAsJsonAsync("api/auth/password-change-request", request);
                return response.IsSuccessStatusCode;
            }
            catch (Exception ex)
            {
                Console.WriteLine($"Password change request error: {ex.Message}");
                return false;
            }
        }

        public async Task<List<PasswordChangeRequest>> GetPasswordChangeRequestsAsync()
        {
            try
            {
                return await _httpClient.GetFromJsonAsync<List<PasswordChangeRequest>>("api/auth/password-change-requests") ?? new List<PasswordChangeRequest>();
            }
            catch (Exception ex)
            {
                Console.WriteLine($"Get password change requests error: {ex.Message}");
                return new List<PasswordChangeRequest>();
            }
        }

        public async Task<bool> ApprovePasswordChangeAsync(PasswordChangeApproval approval)
        {
            try
            {
                var response = await _httpClient.PostAsJsonAsync("api/auth/approve-password-change", approval);
                return response.IsSuccessStatusCode;
            }
            catch (Exception ex)
            {
                Console.WriteLine($"Approve password change error: {ex.Message}");
                return false;
            }
        }

        // Profile Delete Management
        public async Task<bool> RequestProfileDeleteAsync(long userId, string reason)
        {
            try
            {
                var response = await _httpClient.PostAsJsonAsync("api/users/profile-delete-request", new { userId, reason });
                return response.IsSuccessStatusCode;
            }
            catch (Exception ex)
            {
                Console.WriteLine($"Profile delete request error: {ex.Message}");
                return false;
            }
        }

        public async Task<List<ProfileDeleteRequest>> GetProfileDeleteRequestsAsync()
        {
            try
            {
                return await _httpClient.GetFromJsonAsync<List<ProfileDeleteRequest>>("api/users/profile-delete-requests") ?? new List<ProfileDeleteRequest>();
            }
            catch (Exception ex)
            {
                Console.WriteLine($"Get profile delete requests error: {ex.Message}");
                return new List<ProfileDeleteRequest>();
            }
        }

        public async Task<bool> ApproveProfileDeleteAsync(ProfileDeleteApproval approval)
        {
            try
            {
                var response = await _httpClient.PostAsJsonAsync("api/users/approve-profile-delete", approval);
                return response.IsSuccessStatusCode;
            }
            catch (Exception ex)
            {
                Console.WriteLine($"Approve profile delete error: {ex.Message}");
                return false;
            }
        }

        // Email Verification
        public async Task<bool> SendEmailVerificationAsync(string email, UserRole role)
        {
            try
            {
                var response = await _httpClient.PostAsJsonAsync("api/auth/send-verification", new { email, role });
                return response.IsSuccessStatusCode;
            }
            catch (Exception ex)
            {
                Console.WriteLine($"Send email verification error: {ex.Message}");
                return false;
            }
        }

        public async Task<bool> VerifyEmailAsync(string token)
        {
            try
            {
                var response = await _httpClient.PostAsJsonAsync("api/auth/verify-email", new { token });
                return response.IsSuccessStatusCode;
            }
            catch (Exception ex)
            {
                Console.WriteLine($"Verify email error: {ex.Message}");
                return false;
            }
        }

        // Token Management
        public async Task<bool> ValidateTokenAsync(string token)
        {
            try
            {
                var response = await _httpClient.PostAsJsonAsync("api/auth/validate-token", new { token });
                return response.IsSuccessStatusCode;
            }
            catch (Exception ex)
            {
                Console.WriteLine($"Validate token error: {ex.Message}");
                return false;
            }
        }

        public async Task<User?> GetCurrentUserAsync()
        {
            return CurrentUser;
        }

        public async Task LogoutAsync()
        {
            IsAuthenticated = false;
            UserName = null;
            CurrentUserRole = UserRole.Client;
            CurrentUser = null;

            OnAuthenticationChanged?.Invoke(false);

            // Clear stored auth data
            await ClearAuthDataAsync();
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
                    CurrentUserRole = authData.Role;
                    CurrentUser = authData.User;
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

        private async Task SetAuthState(AuthResponse? authResponse)
        {
            if (authResponse != null)
            {
                IsAuthenticated = true;
                UserName = authResponse.User.Username;
                CurrentUserRole = authResponse.User.Role;
                CurrentUser = authResponse.User;

                OnAuthenticationChanged?.Invoke(true);

                // Store in session storage (in real app, use secure tokens)
                await StoreAuthDataAsync(authResponse);
            }
        }

        private async Task StoreAuthDataAsync(AuthResponse authResponse)
        {
            var json = JsonSerializer.Serialize(authResponse);
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
            public UserRole Role { get; set; } = UserRole.Client;
            public User? User { get; set; }
        }
    }
}
