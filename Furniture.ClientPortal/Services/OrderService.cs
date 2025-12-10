using System.Net.Http.Json;
using Furniture.ClientPortal.Models;

namespace Furniture.ClientPortal.Services
{
    public interface IOrderService
    {
        Task<Order?> CreateOrderAsync(CreateOrderRequest request);
        Task<List<Order>> GetOrdersAsync();
        Task<Order?> GetOrderAsync(long id);
        Task SendOrderConfirmationEmailAsync(Order order);
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

        public async Task SendOrderConfirmationEmailAsync(Order order)
        {
            var emailRequest = new SendMailRequest
            {
                To = "test@test.com", 
                Subject = $"Order Confirmation: {order.OrderNumber}",
                Body = $"<h1>Thank you for your order!</h1><p>Your order with number {order.OrderNumber} has been placed successfully.</p>"
            };

            await _httpClient.PostAsJsonAsync("api/mail/send", emailRequest);
        }
    }

    public class SendMailRequest
    {
        public string To { get; set; }
        public string Subject { get; set; }
        public string Body { get; set; }
    }
}
