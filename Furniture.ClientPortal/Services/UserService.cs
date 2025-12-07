using Furniture.Admin.Models;
using System.Security.Cryptography;
using System.Text;
using System.Text.Json;
using Microsoft.AspNetCore.Components;

namespace Furniture.Admin.Services
{
    public interface IUserService
    {
        Task<UserResponse?> CreateUserAsync(CreateUserRequest request);
        Task<UserResponse?> AuthenticateAsync(LoginRequest request);
        Task<UserResponse?> GetUserByIdAsync(long userId);
        Task<UserResponse?> GetUserByUsernameAsync(string username);
        Task<bool> UpdateUserProfileAsync(long userId, UserProfile profile);
        Task<bool> DeactivateUserAsync(long userId);
    }

    public class UserService : IUserService
    {
        private readonly HttpClient _httpClient;
        private readonly NavigationManager _navigationManager;

        public UserService(HttpClient httpClient, NavigationManager navigationManager)
        {
            _httpClient = httpClient;
            _navigationManager = navigationManager;
        }

        public async Task<UserResponse?> CreateUserAsync(CreateUserRequest request)
        {
            try
            {
                // Create user object for API
                // Note: Backend will hash the password with bcrypt
                var user = new
                {
                    username = request.Username,
                    email = request.Email,
                    password = request.Password, // Send plain password - backend will hash it
                    full_name = request.FullName,
                    account_type = request.AccountType,
                    phone = request.Phone,
                    address = request.Address,
                    city = request.City,
                    state = request.State,
                    postal_code = request.PostalCode,
                    country = request.Country,
                    company_name = request.CompanyName
                };

                Console.WriteLine($"Sending request to: {_httpClient.BaseAddress}api/users");
                Console.WriteLine($"Creating user: {request.Username} ({request.Email})");

                // Send to API
                var response = await _httpClient.PostAsJsonAsync("/api/users", user);
                
                Console.WriteLine($"Response status: {response.StatusCode}");
                
                if (response.IsSuccessStatusCode)
                {
                    var createdUser = await response.Content.ReadFromJsonAsync<UserResponse>();
                    return createdUser;
                }
                else
                {
                    var errorContent = await response.Content.ReadAsStringAsync();
                    Console.WriteLine($"API Error: {response.StatusCode} - {errorContent}");
                    return null;
                }
            }
            catch (HttpRequestException ex)
            {
                Console.WriteLine($"HTTP Request Error: {ex.Message}");
                Console.WriteLine("API server may not be running or database not connected");
                return null;
            }
            catch (Exception ex)
            {
                Console.WriteLine($"Error creating user: {ex.Message}");
                throw new Exception($"Registration failed: {ex.Message}");
            }
        }

        public async Task<UserResponse?> AuthenticateAsync(LoginRequest request)
        {
            try
            {
                var loginData = new
                {
                    username = request.Username,
                    password = request.Password
                };

                var response = await _httpClient.PostAsJsonAsync("/api/auth/login", loginData);
                
                if (response.IsSuccessStatusCode)
                {
                    return await response.Content.ReadFromJsonAsync<UserResponse>();
                }

                return null;
            }
            catch (Exception ex)
            {
                Console.WriteLine($"Error authenticating user: {ex.Message}");
                return null;
            }
        }

        public async Task<UserResponse?> GetUserByIdAsync(long userId)
        {
            try
            {
                var response = await _httpClient.GetAsync($"/api/users/{userId}");
                
                if (response.IsSuccessStatusCode)
                {
                    return await response.Content.ReadFromJsonAsync<UserResponse>();
                }

                return null;
            }
            catch (Exception ex)
            {
                Console.WriteLine($"Error getting user: {ex.Message}");
                return null;
            }
        }

        public async Task<UserResponse?> GetUserByUsernameAsync(string username)
        {
            try
            {
                var response = await _httpClient.GetAsync($"/api/users/username/{username}");
                
                if (response.IsSuccessStatusCode)
                {
                    return await response.Content.ReadFromJsonAsync<UserResponse>();
                }

                return null;
            }
            catch (Exception ex)
            {
                Console.WriteLine($"Error getting user by username: {ex.Message}");
                return null;
            }
        }

        public async Task<bool> UpdateUserProfileAsync(long userId, UserProfile profile)
        {
            try
            {
                var response = await _httpClient.PutAsJsonAsync($"/api/users/{userId}/profile", profile);
                return response.IsSuccessStatusCode;
            }
            catch (Exception ex)
            {
                Console.WriteLine($"Error updating user profile: {ex.Message}");
                return false;
            }
        }

        public async Task<bool> DeactivateUserAsync(long userId)
        {
            try
            {
                var response = await _httpClient.DeleteAsync($"/api/users/{userId}");
                return response.IsSuccessStatusCode;
            }
            catch (Exception ex)
            {
                Console.WriteLine($"Error deactivating user: {ex.Message}");
                return false;
            }
        }

        private string HashPassword(string password)
        {
            using var sha256 = SHA256.Create();
            var hashedBytes = sha256.ComputeHash(Encoding.UTF8.GetBytes(password));
            return Convert.ToBase64String(hashedBytes);
        }

        private bool VerifyPassword(string password, string hash)
        {
            var computedHash = HashPassword(password);
            return computedHash == hash;
        }
    }
}
