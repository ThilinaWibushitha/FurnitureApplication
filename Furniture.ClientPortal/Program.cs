using Microsoft.AspNetCore.Components;
using Microsoft.AspNetCore.Components.Web;

var builder = WebApplication.CreateBuilder(args);

// Add services
builder.Services.AddRazorPages();
builder.Services.AddServerSideBlazor();
builder.Services.AddScoped(sp => new HttpClient { BaseAddress = new Uri("http://localhost:8081/") });
builder.Services.AddScoped<Furniture.Admin.Services.IItemService, Furniture.Admin.Services.ItemService>();
builder.Services.AddScoped<Furniture.Admin.Services.AuthenticationService>();
builder.Services.AddScoped<Furniture.Admin.Services.IUserService, Furniture.Admin.Services.UserService>();
builder.Services.AddScoped<Furniture.Admin.Services.IPaymentService, Furniture.Admin.Services.PaymentService>();
builder.Services.AddScoped<Furniture.Admin.Services.IRatingService, Furniture.Admin.Services.RatingService>();
builder.Services.AddScoped<Furniture.Admin.Services.IProfileService, Furniture.Admin.Services.ProfileService>();


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
