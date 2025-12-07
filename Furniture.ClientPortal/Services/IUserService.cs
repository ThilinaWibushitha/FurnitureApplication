namespace Furniture.ClientPortal.Services;

public interface IUserService
{
    Task CreateAccount(string email, string password);
}
