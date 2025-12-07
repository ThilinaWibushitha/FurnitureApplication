using System.Net.Http.Json;
using Furniture.ClientPortal.Models;

namespace Furniture.ClientPortal.Services
{
    public interface IOrderService
    {
        Task<Order?> CreateOrderAsync(CreateOrderRequest request);
        Task<List<Order>> GetOrdersAsync();
        Task<Order?> GetOrderAsync(long id);
    }

    public class OrderService : IOrderService
    {
        private readonly HttpClient _httpClient;

        public OrderService(HttpClient httpClient)
        {
            _httpClient = httpClient;
        }

        public async Task<Order?> CreateOrderAsync(CreateOrderRequest request)
        {
            var response = await _httpClient.PostAsJsonAsync("api/orders", request);
            if (response.IsSuccessStatusCode)
            {
                return await response.Content.ReadFromJsonAsync<Order>();
            }
            return null;
        }

        public async Task<List<Order>> GetOrdersAsync()
        {
             return await _httpClient.GetFromJsonAsync<List<Order>>("api/orders") ?? new List<Order>();
        }

        public async Task<Order?> GetOrderAsync(long id)
        {
             return await _httpClient.GetFromJsonAsync<Order>($"api/orders/{id}");
        }
    }
}
