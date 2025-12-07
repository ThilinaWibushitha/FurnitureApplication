using Furniture.ClientPortal.Models;

namespace Furniture.ClientPortal.Services
{
    public interface IAuthenticationService
    {
        // Main Admin Authentication
        Task<AuthResponse?> MainAdminLoginAsync(MainAdminLoginRequest request);
        Task<bool> MainAdminPasswordResetRequestAsync(string email);
        Task<bool> ConfirmPasswordResetAsync(PasswordResetConfirm confirm);

        // Regular Admin Authentication
        Task<AuthResponse?> AdminLoginAsync(AdminLoginRequest request);
        Task<bool> AdminPasswordResetRequestAsync(string username);

        // Client Authentication
        Task<AuthResponse?> ClientLoginAsync(LoginRequest request);
        Task<bool> ClientPasswordResetRequestAsync(string email);

        // User Management (Main Admin only)
        Task<User?> CreateUserAsync(CreateUserRequest request);
        Task<User?> CreateAdminAsync(CreateAdminRequest request);
        Task<bool> DeleteUserAsync(long userId);
        Task<List<User>> GetUsersAsync();

        // Password Change Management
        Task<bool> RequestPasswordChangeAsync(PasswordChangeRequest request);
        Task<List<PasswordChangeRequest>> GetPasswordChangeRequestsAsync();
        Task<bool> ApprovePasswordChangeAsync(PasswordChangeApproval approval);

        // Profile Delete Management
        Task<bool> RequestProfileDeleteAsync(long userId, string reason);
        Task<List<ProfileDeleteRequest>> GetProfileDeleteRequestsAsync();
        Task<bool> ApproveProfileDeleteAsync(ProfileDeleteApproval approval);

        // Email Verification
        Task<bool> SendEmailVerificationAsync(string email, UserRole role);
        Task<bool> VerifyEmailAsync(string token);

        // Token Management
        Task<bool> ValidateTokenAsync(string token);
        Task<User?> GetCurrentUserAsync();
        Task LogoutAsync();

        // Events
        event Action<bool>? OnAuthenticationChanged;

        // Auth Status
        Task<bool> CheckAuthStatusAsync();
    }
}
