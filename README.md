# PlazaNet Account Server
A Go Gin Server providing accounts for [PlazaNet](https://github.com/PlazaNetOrg/PlazaNet) and [GamePlaza](https://github.com/PlazaNetOrg/GamePlaza)

## Features:
- **Authentication & Session Management**:
  - JWT tokens (24h) + cookie-based sessions for web
  - Middleware supports both Bearer token and cookie

- **Security & Configuration**:
  - bcrypt password hashing
  - JWT signing via `.env` (`JWT_SECRET`)
  - Configurable cookie flags (`AUTH_SECURE`, `HTTP_ONLY`)
  - SQLite + GORM (auto-migrates `accounts.db`)

- **In progress / planned**:
  - Pals
    - Previously 2D avatars that worked in the Python version. Currently being reimagined to 3D. (Three.js)
  - User status (online/offline/playing)
  - Friends system (add/remove/list)
  - Account editing and deleting.

## Self-Hosted Setup:
You can find more info about self-hosting the Account Server in [PlazaNet Docs](https://plazanet.org/docs/#accounts-install)

## Support:
Wanna support the project? Here are the best ways you can:
- Contribute to the code:
  - Implement new features
  - Fix bugs
- Report bugs you can't fix yourself
- Donate to Andus on [Ko-fi](https://ko-fi.com/andusdev)
  - This will help with payments for the domain and VPS hosting for the official PlazaNet Instance.
> [!IMPORTANT]  
> This project is still in development and isn’t ready to be hosted publicly yet. If you want to support my work in general, feel free to donate, but if your goal is to directly support PlazaNet, please wait until it’s ready for public hosting.

## Credits:
- Andus - Developer, Artist
