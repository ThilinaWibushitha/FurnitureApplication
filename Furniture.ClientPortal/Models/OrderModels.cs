namespace Furniture.ClientPortal.Models
{
    public class Order
    {
        public long Id { get; set; }
        public string OrderNumber { get; set; } = string.Empty;
        public long? CustomerId { get; set; }
        public long? UserId { get; set; }
        public long? BranchId { get; set; }
        public string Status { get; set; } = "Pending";
        public decimal TotalAmount { get; set; }
        public decimal TaxAmount { get; set; }
        public decimal DiscountAmount { get; set; }
        public string PaymentMethod { get; set; } = "CreditCard";
        public string Notes { get; set; } = string.Empty;
        public DateTime CreatedAt { get; set; }
        public DateTime UpdatedAt { get; set; }
    }

    public class OrderItem
    {
        public long Id { get; set; }
        public long OrderId { get; set; }
        public long ItemVariantId { get; set; } // We might use ItemId for now if no variants
        public int Quantity { get; set; }
        public decimal UnitPrice { get; set; }
        public decimal Subtotal { get; set; }
        public decimal Discount { get; set; }
    }

    public class CreateOrderRequest
    {
        // Flattened Order properties
        public string OrderNumber { get; set; } = string.Empty;
        public long? CustomerId { get; set; }
        public long? UserId { get; set; }
        public long? BranchId { get; set; }
        public string Status { get; set; } = "Pending";
        public decimal TotalAmount { get; set; }
        public decimal TaxAmount { get; set; }
        public decimal DiscountAmount { get; set; }
        public string PaymentMethod { get; set; } = "CreditCard";
        public string Notes { get; set; } = string.Empty;

        public List<OrderItem> Items { get; set; } = new List<OrderItem>();
    }
}
