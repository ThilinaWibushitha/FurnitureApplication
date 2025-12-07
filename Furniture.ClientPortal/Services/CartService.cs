using Furniture.ClientPortal.Models;

namespace Furniture.ClientPortal.Services
{
    public class CartItem
    {
        public Item Item { get; set; } = new();
        public int Quantity { get; set; }
        public decimal SubTotal => (decimal)Item.BasePrice * Quantity;
    }

    public interface ICartService
    {
        event Action? OnChange;
        List<CartItem> Items { get; }
        void AddToCart(Item item, int quantity = 1);
        void RemoveFromCart(long itemId);
        void UpdateQuantity(long itemId, int quantity);
        void ClearCart();
        decimal TotalPrice();
        int TotalCount();
    }

    public class CartService : ICartService
    {
        public event Action? OnChange;
        public List<CartItem> Items { get; private set; } = new List<CartItem>();

        public void AddToCart(Item item, int quantity = 1)
        {
            var cartItem = Items.FirstOrDefault(x => x.Item.Id == item.Id);
            if (cartItem == null)
            {
                Items.Add(new CartItem { Item = item, Quantity = quantity });
            }
            else
            {
                cartItem.Quantity += quantity;
            }
            NotifyStateChanged();
        }

        public void RemoveFromCart(long itemId)
        {
            var cartItem = Items.FirstOrDefault(x => x.Item.Id == itemId);
            if (cartItem != null)
            {
                Items.Remove(cartItem);
                NotifyStateChanged();
            }
        }

        public void UpdateQuantity(long itemId, int quantity)
        {
            var cartItem = Items.FirstOrDefault(x => x.Item.Id == itemId);
            if (cartItem != null)
            {
                if (quantity <= 0)
                {
                    Items.Remove(cartItem);
                }
                else
                {
                    cartItem.Quantity = quantity;
                }
                NotifyStateChanged();
            }
        }

        public void ClearCart()
        {
            Items.Clear();
            NotifyStateChanged();
        }

        public decimal TotalPrice()
        {
            return Items.Sum(x => x.SubTotal);
        }

        public int TotalCount()
        {
            return Items.Sum(x => x.Quantity);
        }

        private void NotifyStateChanged() => OnChange?.Invoke();
    }
}
