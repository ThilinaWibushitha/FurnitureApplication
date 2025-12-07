using System.Net.Http.Json;
using Furniture.ClientPortal.Models;

namespace Furniture.ClientPortal.Services
{
    public interface IPaymentService
    {
        Task<List<Payment>> GetPaymentsAsync(string? status = null);
        Task<Payment?> GetPaymentAsync(long id);
        Task<List<Payment>> GetOrderPaymentsAsync(long orderId);
        Task<Payment> CreatePaymentAsync(CreatePaymentRequest request);
        Task CancelPaymentAsync(long id, string reason);
    }

    public class PaymentService : IPaymentService
    {
        private readonly HttpClient _httpClient;

        public PaymentService(HttpClient httpClient)
        {
            _httpClient = httpClient;
        }

        public async Task<List<Payment>> GetPaymentsAsync(string? status = null)
        {
            var url = "api/payments";
            if (!string.IsNullOrEmpty(status))
                url += $"?status={status}";
            
            return await _httpClient.GetFromJsonAsync<List<Payment>>(url) ?? new List<Payment>();
        }

        public async Task<Payment?> GetPaymentAsync(long id)
        {
            return await _httpClient.GetFromJsonAsync<Payment>($"api/payments/{id}");
        }

        public async Task<List<Payment>> GetOrderPaymentsAsync(long orderId)
        {
            return await _httpClient.GetFromJsonAsync<List<Payment>>($"api/orders/{orderId}/payments") ?? new List<Payment>();
        }

        public async Task<Payment> CreatePaymentAsync(CreatePaymentRequest request)
        {
            var response = await _httpClient.PostAsJsonAsync("api/payments", request);
            response.EnsureSuccessStatusCode();
            return await response.Content.ReadFromJsonAsync<Payment>() ?? throw new Exception("Failed to create payment");
        }

        public async Task CancelPaymentAsync(long id, string reason)
        {
            await _httpClient.PostAsJsonAsync($"api/payments/{id}/cancel", new { reason });
        }
    }
}
