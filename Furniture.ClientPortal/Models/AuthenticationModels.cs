using System.ComponentModel.DataAnnotations;

namespace Furniture.ClientPortal.Models
{
    public class LoginRequest
    {
        [Required]
        public string Username { get; set; } = string.Empty;

        [Required]
        public string Password { get; set; } = string.Empty;

        public UserRole Role { get; set; } = UserRole.Client;
    }

    public class MainAdminLoginRequest
    {
        [Required]
        [EmailAddress]
        public string Email { get; set; } = string.Empty;

        [Required]
        public string Password { get; set; } = string.Empty;
    }

    public class AdminLoginRequest
    {
        [Required]
        public string Username { get; set; } = string.Empty;

        [Required]
        public string Password { get; set; } = string.Empty;
    }

    public class PasswordResetRequest
    {
        [Required]
        [EmailAddress]
        public string Email { get; set; } = string.Empty;

        public UserRole Role { get; set; } = UserRole.Client;
    }

    public class PasswordResetConfirm
    {
        [Required]
        public string Token { get; set; } = string.Empty;

        [Required]
        [StringLength(100, MinimumLength = 6)]
        public string NewPassword { get; set; } = string.Empty;

        [Required]
        [Compare(nameof(NewPassword), ErrorMessage = "Passwords do not match")]
        public string ConfirmPassword { get; set; } = string.Empty;
    }

    public class PasswordChangeRequest
    {
        [Required]
        public long UserId { get; set; }

        [Required]
        [StringLength(100, MinimumLength = 6)]
        public string NewPassword { get; set; } = string.Empty;

        [Required]
        public string Reason { get; set; } = string.Empty;

        public bool RequiresApproval { get; set; } = true;
    }

    public class PasswordChangeApproval
    {
        [Required]
        public long RequestId { get; set; }

        [Required]
        public bool Approved { get; set; }

        public string? AdminNotes { get; set; }
    }

    public class ProfileDeleteRequest
    {
        [Required]
        public long UserId { get; set; }

        [Required]
        public string Reason { get; set; } = string.Empty;

        public DateTime RequestedAt { get; set; } = DateTime.UtcNow;

        public bool IsApproved { get; set; } = false;

        public DateTime? ApprovedAt { get; set; }

        public long? ApprovedBy { get; set; }

        public string? AdminNotes { get; set; }
    }

    public class ProfileDeleteApproval
    {
        [Required]
        public long RequestId { get; set; }

        [Required]
        public bool Approved { get; set; }

        public string? AdminNotes { get; set; }
    }

    public class CreateAdminRequest
    {
        [Required]
        [StringLength(50, MinimumLength = 3)]
        public string Username { get; set; } = string.Empty;

        [Required]
        [EmailAddress]
        public string Email { get; set; } = string.Empty;

        [Required]
        [StringLength(100, MinimumLength = 6)]
        public string Password { get; set; } = string.Empty;

        [Required]
        [Compare(nameof(Password), ErrorMessage = "Passwords do not match")]
        public string ConfirmPassword { get; set; } = string.Empty;

        [Required]
        [StringLength(100)]
        public string FullName { get; set; } = string.Empty;

        public UserRole Role { get; set; } = UserRole.Admin; // Main admin creates regular admins
    }

    public class AuthResponse
    {
        public string Token { get; set; } = string.Empty;
        public User User { get; set; } = null!;
        public DateTime ExpiresAt { get; set; }
        public string[] Permissions { get; set; } = Array.Empty<string>();
    }

    public class EmailVerificationRequest
    {
        [Required]
        [EmailAddress]
        public string Email { get; set; } = string.Empty;

        public UserRole Role { get; set; } = UserRole.Client;

        public string VerificationToken { get; set; } = string.Empty;

        public DateTime ExpiresAt { get; set; } = DateTime.UtcNow.AddHours(24);

        public bool IsVerified { get; set; } = false;

        public DateTime? VerifiedAt { get; set; }
    }
}
