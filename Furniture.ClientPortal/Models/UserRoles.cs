namespace Furniture.ClientPortal.Models
{
    public enum UserRole
    {
        MainAdmin = 0,
        Admin = 1,
        Client = 2
    }

    public static class UserRoleExtensions
    {
        public static string GetDisplayName(this UserRole role)
        {
            return role switch
            {
                UserRole.MainAdmin => "Main Administrator",
                UserRole.Admin => "Administrator",
                UserRole.Client => "Client",
                _ => "Unknown"
            };
        }

        public static bool CanAccessMainAdminPanel(this UserRole role)
        {
            return role == UserRole.MainAdmin;
        }

        public static bool CanAccessAdminPanel(this UserRole role)
        {
            return role == UserRole.MainAdmin || role == UserRole.Admin;
        }

        public static bool CanAccessClientPanel(this UserRole role)
        {
            return role == UserRole.Client;
        }

        public static bool CanManageUsers(this UserRole role)
        {
            return role == UserRole.MainAdmin;
        }

        public static bool CanManageItems(this UserRole role)
        {
            return role == UserRole.MainAdmin || role == UserRole.Admin;
        }

        public static bool CanViewSales(this UserRole role)
        {
            return role == UserRole.MainAdmin;
        }

        public static bool CanManageOrders(this UserRole role)
        {
            return role == UserRole.MainAdmin || role == UserRole.Admin;
        }

        public static bool CanManageClientProfiles(this UserRole role)
        {
            return role == UserRole.MainAdmin || role == UserRole.Admin;
        }

        public static bool CanApproveDeleteRequests(this UserRole role)
        {
            return role == UserRole.MainAdmin || role == UserRole.Admin;
        }

        public static bool CanApprovePasswordChanges(this UserRole role)
        {
            return role == UserRole.MainAdmin;
        }

        public static bool CanSendEmails(this UserRole role)
        {
            return role == UserRole.MainAdmin || role == UserRole.Admin;
        }
    }
}
