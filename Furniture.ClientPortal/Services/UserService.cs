using Furniture.ClientPortal.Models;

namespace Furniture.ClientPortal.Services;

public class UserService : IUserService
{
    public Task CreateAccount(string email, string password)
    {
        // In a real application, you would save the user to a database.
        // For this example, we'll just simulate a delay.
        return Task.Delay(1000);
    }
}
