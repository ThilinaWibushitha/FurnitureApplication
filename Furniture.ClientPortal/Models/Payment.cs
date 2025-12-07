namespace Furniture.Admin.Models
{
    public class Payment
    {
        public long Id { get; set; }
        public long OrderId { get; set; }
        public long? CustomerId { get; set; }
        public decimal Amount { get; set; }
        public string PaymentMethod { get; set; } = "Cash";
        public string PaymentStatus { get; set; } = "Pending";
        public string? TransactionRef { get; set; }
        public string? Notes { get; set; }
        public DateTime CreatedAt { get; set; }
        public DateTime UpdatedAt { get; set; }
    }

    public class CreatePaymentRequest
    {
        public long OrderId { get; set; }
        public decimal Amount { get; set; }
        public string PaymentMethod { get; set; } = "Cash";
        public string? Notes { get; set; }
    }

    public class ItemRating
    {
        public long Id { get; set; }
        public long ItemId { get; set; }
        public long CustomerId { get; set; }
        public long? OrderId { get; set; }
        public int Rating { get; set; }
        public string? ReviewText { get; set; }
        public bool IsVerifiedPurchase { get; set; }
        public bool IsApproved { get; set; }
        public DateTime CreatedAt { get; set; }
        public string? CustomerName { get; set; }
    }

    public class CreateRatingRequest
    {
        public long ItemId { get; set; }
        public int Rating { get; set; }
        public string? ReviewText { get; set; }
        public long? OrderId { get; set; }
    }

    public class RatingSummary
    {
        public long ItemId { get; set; }
        public int TotalRatings { get; set; }
        public double AverageRating { get; set; }
        public int Rating5Count { get; set; }
        public int Rating4Count { get; set; }
        public int Rating3Count { get; set; }
        public int Rating2Count { get; set; }
        public int Rating1Count { get; set; }
    }

    public class ClientProfile
    {
        public long CustomerId { get; set; }
        public string FirstName { get; set; } = "";
        public string? LastName { get; set; }
        public string? Email { get; set; }
        public string? Phone { get; set; }
        public string? Address { get; set; }
        public string? City { get; set; }
        public int LoyaltyPoints { get; set; }
        public string? DateOfBirth { get; set; }
        public string? Gender { get; set; }
        public bool NewsletterOptIn { get; set; }
        public string? ProfileImageUrl { get; set; }
        public string AccountStatus { get; set; } = "Active";
    }

    public class UpdateProfileRequest
    {
        public string? Phone { get; set; }
        public string? Address { get; set; }
        public string? City { get; set; }
        public string? DateOfBirth { get; set; }
        public string? Gender { get; set; }
        public bool? NewsletterOptIn { get; set; }
        public string? ProfileImageUrl { get; set; }
    }

    public class Order
    {
        public long Id { get; set; }
        public string OrderNumber { get; set; } = "";
        public long? CustomerId { get; set; }
        public long? UserId { get; set; }
        public long? BranchId { get; set; }
        public string Status { get; set; } = "Pending";
        public decimal TotalAmount { get; set; }
        public decimal TaxAmount { get; set; }
        public decimal DiscountAmount { get; set; }
        public string? PaymentMethod { get; set; }
        public string? Notes { get; set; }
        public DateTime CreatedAt { get; set; }
        public DateTime UpdatedAt { get; set; }
        public List<OrderItem> Items { get; set; } = new();
    }

    public class OrderItem
    {
        public long Id { get; set; }
        public long OrderId { get; set; }
        public long? ItemVariantId { get; set; }
        public int Quantity { get; set; }
        public decimal UnitPrice { get; set; }
        public decimal Subtotal { get; set; }
        public decimal Discount { get; set; }
    }
}
