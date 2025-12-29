using Microsoft.AspNetCore.Components;
using Microsoft.AspNetCore.Components.Web;

var builder = WebApplication.CreateBuilder(args);

// Add services
builder.Services.AddRazorPages();
builder.Services.AddServerSideBlazor();
builder.Services.AddScoped(sp => new HttpClient { BaseAddress = new Uri("http://localhost:8080/") });
builder.Services.AddScoped<Furniture.ClientPortal.Services.IAuthenticationService, Furniture.ClientPortal.Services.AuthenticationService>();
builder.Services.AddScoped<Furniture.ClientPortal.Services.IUserService, Furniture.ClientPortal.Services.UserService>();
builder.Services.AddScoped<Furniture.ClientPortal.Services.ICartService, Furniture.ClientPortal.Services.CartService>();
builder.Services.AddScoped<Furniture.ClientPortal.Services.IItemService, Furniture.ClientPortal.Services.ItemService>();
builder.Services.AddScoped<Furniture.ClientPortal.Services.IPaymentService, Furniture.ClientPortal.Services.PaymentService>();
builder.Services.AddScoped<Furniture.ClientPortal.Services.IRatingService, Furniture.ClientPortal.Services.RatingService>();
builder.Services.AddScoped<Furniture.ClientPortal.Services.IProfileService, Furniture.ClientPortal.Services.ProfileService>();


var app = builder.Build();

// Configure the HTTP request pipeline.
if (!app.Environment.IsDevelopment()) {
    app.UseExceptionHandler("/Error");
    app.UseHsts();
}

app.UseStaticFiles();
app.UseRouting();

app.MapBlazorHub();
app.MapRazorPages();
app.MapFallbackToPage("/_Host");

app.Run();
