using System.Text.Json.Serialization;

namespace Furniture.Admin.Models
{
    public class Item
    {
        [JsonPropertyName("id")]
        public long Id { get; set; }

        [JsonPropertyName("sku")]
        public string Sku { get; set; } = string.Empty;

        [JsonPropertyName("name")]
        public string Name { get; set; } = string.Empty;

        [JsonPropertyName("description")]
        public string Description { get; set; } = string.Empty;

        [JsonPropertyName("department_id")]
        public long DepartmentId { get; set; }

        [JsonPropertyName("category_id")]
        public long CategoryId { get; set; }

        [JsonPropertyName("main_image_url")]
        public string MainImageUrl { get; set; } = string.Empty;

        [JsonPropertyName("status")]
        public string Status { get; set; } = "Active";

        [JsonPropertyName("base_price")]
        public double BasePrice { get; set; }

        [JsonPropertyName("base_cost")]
        public double BaseCost { get; set; }

        [JsonPropertyName("default_discount_percent")]
        public double DefaultDiscountPercent { get; set; }

        [JsonPropertyName("tax_class_id")]
        public long TaxClassId { get; set; }

        [JsonPropertyName("is_published_online")]
        public bool IsPublishedOnline { get; set; }

        [JsonPropertyName("online_sort_order")]
        public int OnlineSortOrder { get; set; }
    }
}
