# PlazaNet Account Server
A Python Flask Server providing accounts for [PlazaNet](https://github.com/PlazaNetOrg/PlazaNet) and [GamePlaza](https://github.com/PlazaNetOrg/GamePlaza)

## Features:
- Basic API Endpoints:
  - Login with JWT tokens for PlazaNet and GamePlaza
  - Get user info
  - Set and get user status
  - Add/remove friends and fetch friends list
  - Return user's Pal
- Account Managment:
  - Account creation and deletion
  - Birthday management
- Avatar (Pal) Editor:
  - Create and customize your Pal
- Security:
  - JWT-based authentication
  - Environment-based configuration (.env) for secrets and settings

## Self-Hosted Setup:
### Linux:
1. Create a directory (if you haven't already): `mkdir plazanet && cd plazanet`
2. Clone the repository using: `git clone https://github.com/PlazaNetOrg/Accounts accounts && cd accounts`
3. Install all requirements:
    - python with pip (`sudo apt install python3 python3-pip`)
    - Python requirements:
      - `pip install -r requirements.txt` 
    - openssl (`sudo apt install openssl` on Ubuntu) for generation of secret keys.
4. Rename `.env.example` to `.env`
5. Replace the values in `.env`:
    - Generate `SECRET_KEY` using `openssl rand -base64 32`
    - Generate `JWT_SECRET_KEY` using `openssl rand -base64 32`
    - Optionally replace the rest of the values.
6. Initialize the database:
    - `flask --app run.py db init`
    - `flask --app run.py db migrate -m "Initial migration"`
    - `flask --app run.py db upgrade`
7. Start the server using `flask --app run.py`

## Credits:
- Andus - Developer, Artist