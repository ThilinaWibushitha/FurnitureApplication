namespace Furniture.Clientportal.Services;

public interface IUserService
{
    Task CreateAccount(string email, string password);
}
