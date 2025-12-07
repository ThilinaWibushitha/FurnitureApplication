using System.Net.Http.Json;
using Furniture.Admin.Models;

namespace Furniture.Admin.Services
{
    public interface IProfileService
    {
        Task<ClientProfile?> GetProfileAsync();
        Task UpdateProfileAsync(UpdateProfileRequest request);
        Task ChangePasswordAsync(string currentPassword, string newPassword);
        Task RequestDeleteAccountAsync(string password, string? reason);
    }

    public class ProfileService : IProfileService
    {
        private readonly HttpClient _httpClient;

        public ProfileService(HttpClient httpClient)
        {
            _httpClient = httpClient;
        }

        public async Task<ClientProfile?> GetProfileAsync()
        {
            try
            {
                return await _httpClient.GetFromJsonAsync<ClientProfile>("api/profile");
            }
            catch
            {
                return null;
            }
        }

        public async Task UpdateProfileAsync(UpdateProfileRequest request)
        {
            var response = await _httpClient.PutAsJsonAsync("api/profile", request);
            response.EnsureSuccessStatusCode();
        }

        public async Task ChangePasswordAsync(string currentPassword, string newPassword)
        {
            var response = await _httpClient.PostAsJsonAsync("api/profile/password", new 
            { 
                current_password = currentPassword, 
                new_password = newPassword 
            });
            response.EnsureSuccessStatusCode();
        }

        public async Task RequestDeleteAccountAsync(string password, string? reason)
        {
            var response = await _httpClient.PostAsJsonAsync("api/profile/delete", new 
            { 
                password, 
                reason 
            });
            response.EnsureSuccessStatusCode();
        }
    }
}
