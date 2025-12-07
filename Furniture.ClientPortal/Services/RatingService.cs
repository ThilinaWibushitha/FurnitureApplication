using System.Net.Http.Json;
using Furniture.Admin.Models;

namespace Furniture.Admin.Services
{
    public interface IRatingService
    {
        Task<List<ItemRating>> GetItemRatingsAsync(long itemId);
        Task<RatingSummary?> GetRatingSummaryAsync(long itemId);
        Task<List<ItemRating>> GetMyRatingsAsync();
        Task<ItemRating> CreateRatingAsync(long itemId, CreateRatingRequest request);
        Task UpdateRatingAsync(long id, int rating, string? reviewText);
        Task DeleteRatingAsync(long id);
    }

    public class RatingService : IRatingService
    {
        private readonly HttpClient _httpClient;

        public RatingService(HttpClient httpClient)
        {
            _httpClient = httpClient;
        }

        public async Task<List<ItemRating>> GetItemRatingsAsync(long itemId)
        {
            return await _httpClient.GetFromJsonAsync<List<ItemRating>>($"api/items/{itemId}/ratings") ?? new List<ItemRating>();
        }

        public async Task<RatingSummary?> GetRatingSummaryAsync(long itemId)
        {
            return await _httpClient.GetFromJsonAsync<RatingSummary>($"api/items/{itemId}/rating-summary");
        }

        public async Task<List<ItemRating>> GetMyRatingsAsync()
        {
            return await _httpClient.GetFromJsonAsync<List<ItemRating>>("api/profile/ratings") ?? new List<ItemRating>();
        }

        public async Task<ItemRating> CreateRatingAsync(long itemId, CreateRatingRequest request)
        {
            var response = await _httpClient.PostAsJsonAsync($"api/items/{itemId}/ratings", request);
            response.EnsureSuccessStatusCode();
            return await response.Content.ReadFromJsonAsync<ItemRating>() ?? throw new Exception("Failed to create rating");
        }

        public async Task UpdateRatingAsync(long id, int rating, string? reviewText)
        {
            await _httpClient.PutAsJsonAsync($"api/ratings/{id}", new { rating, review_text = reviewText });
        }

        public async Task DeleteRatingAsync(long id)
        {
            await _httpClient.DeleteAsync($"api/ratings/{id}");
        }
    }
}
