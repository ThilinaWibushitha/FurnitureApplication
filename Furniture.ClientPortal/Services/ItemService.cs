using System.Net.Http.Json;
using Furniture.ClientPortal.Models;

namespace Furniture.ClientPortal.Services
{
    public interface IItemService
    {
        Task<List<Item>> GetItemsAsync();
        Task<Item?> GetItemAsync(long id);
        Task CreateItemAsync(Item item);
        Task UpdateItemAsync(Item item);
        Task DeleteItemAsync(long id);
    }

    public class ItemService : IItemService
    {
        private readonly HttpClient _httpClient;

        public ItemService(HttpClient httpClient)
        {
            _httpClient = httpClient;
        }

        public async Task<List<Item>> GetItemsAsync()
        {
            return await _httpClient.GetFromJsonAsync<List<Item>>("api/items") ?? new List<Item>();
        }

        public async Task<Item?> GetItemAsync(long id)
        {
            return await _httpClient.GetFromJsonAsync<Item>($"api/items/{id}");
        }

        public async Task CreateItemAsync(Item item)
        {
            await _httpClient.PostAsJsonAsync("api/items", item);
        }

        public async Task UpdateItemAsync(Item item)
        {
            await _httpClient.PutAsJsonAsync($"api/items/{item.Id}", item);
        }

        public async Task DeleteItemAsync(long id)
        {
            await _httpClient.DeleteAsync($"api/items/{id}");
        }
    }
}
