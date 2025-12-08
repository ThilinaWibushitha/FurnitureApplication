using System;
using System.ComponentModel.DataAnnotations;
using System.Text.Json.Serialization;

namespace Furniture.ClientPortal.Models
{
    public class User
    {
        public long Id { get; set; }
        
        [Required]
        [StringLength(50, MinimumLength = 3)]
        public string Username { get; set; } = string.Empty;
        
        [Required]
        [EmailAddress]
        public string Email { get; set; } = string.Empty;
        
        [Required]
        public string PasswordHash { get; set; } = string.Empty;
        
        [Required]
        [StringLength(100)]
        public string FullName { get; set; } = string.Empty;
        
        [Required]
        public UserRole Role { get; set; } = UserRole.Client;
        
        public bool IsActive { get; set; } = true;
        public bool IsVerified { get; set; } = false;
        public DateTime CreatedAt { get; set; }
        public DateTime UpdatedAt { get; set; }
        public DateTime? LastLoginAt { get; set; }
        
        // Navigation property
        public UserProfile? Profile { get; set; }
    }

    public class UserProfile
    {
        public long Id { get; set; }
        public long UserId { get; set; }
        
        [Phone]
        public string? Phone { get; set; }
        
        public string? Address { get; set; }
        public string? City { get; set; }
        public string? State { get; set; }
        public string? PostalCode { get; set; }
        public string Country { get; set; } = "US";
        public string? CompanyName { get; set; }
        public string? ProfileImageUrl { get; set; }
        public string Preferences { get; set; } = "{}";
        public DateTime CreatedAt { get; set; }
        public DateTime UpdatedAt { get; set; }
        
        // Navigation property
        public User User { get; set; } = null!;
    }

    public class CreateUserRequest
    {
        [Required]
        [StringLength(100)]
        [JsonPropertyName("full_name")]
        public string FullName { get; set; } = string.Empty;

        [Required]
        [EmailAddress]
        [JsonPropertyName("email")]
        public string Email { get; set; } = string.Empty;

        [Required]
        [StringLength(50, MinimumLength = 3)]
        [JsonPropertyName("username")]
        public string Username { get; set; } = string.Empty;

        [Required]
        [StringLength(100, MinimumLength = 6)]
        [JsonPropertyName("password")]
        public string Password { get; set; } = string.Empty;

        [Required]
        [Compare(nameof(Password), ErrorMessage = "Passwords do not match")]
        [JsonPropertyName("confirm_password")]
        public string ConfirmPassword { get; set; } = string.Empty;

        [JsonIgnore]
        public UserRole Role { get; set; } = UserRole.Client;

        [JsonPropertyName("account_type")]
        public string AccountType
        {
            get => Role switch
            {
                UserRole.MainAdmin => "admin",
                UserRole.Admin => "admin",
                _ => "client"
            };
            set
            {
                if (string.Equals(value, "admin", StringComparison.OrdinalIgnoreCase))
                {
                    Role = UserRole.Admin;
                }
                else
                {
                    Role = UserRole.Client;
                }
            }
        }

        [Range(typeof(bool), "true", "true", ErrorMessage = "You must agree to the terms")]
        [JsonPropertyName("agree_terms")]
        public bool AgreeTerms { get; set; }

        [Phone]
        [JsonPropertyName("phone")]
        public string? Phone { get; set; }

        [JsonPropertyName("address")]
        public string? Address { get; set; }

        [JsonPropertyName("city")]
        public string? City { get; set; }

        [JsonPropertyName("state")]
        public string? State { get; set; }

        [JsonPropertyName("postal_code")]
        public string? PostalCode { get; set; }

        [JsonPropertyName("country")]
        public string Country { get; set; } = "US";

        [JsonPropertyName("company_name")]
        public string? CompanyName { get; set; }
    }

    public class UserResponse
    {
        public long Id { get; set; }
        public string Username { get; set; } = string.Empty;
        public string Email { get; set; } = string.Empty;
        public string FullName { get; set; } = string.Empty;
        public UserRole Role { get; set; } = UserRole.Client;
        public bool IsActive { get; set; }
        public bool IsVerified { get; set; }
        public DateTime CreatedAt { get; set; }
        public DateTime? LastLoginAt { get; set; }
        public UserProfile? Profile { get; set; }
    }
}
