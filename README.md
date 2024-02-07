# Overview
This project shows a very bare bones example of using Azure Entra ID as an Authorization Source with a golang API
protected Endpoint. This was created to educate myself, and hopefully others on how to properly configure everything
needed to get a working application setup.

# Steps
This section will go in-depth on what is needed to get up and running with a protected API endpoint. Since we have the example
code handy, this section mainly will be covering everything Needed from Azure's standpoint. I will not be walking through
how to get started with azure, how to create an account, or anything of that nature. If you need that information, please
refer to google or chatGPT.

## Azure
### Creating an App Registration
We are going to create an Azure App Reg. 
1. Log into Azure and navigate to Microsoft Entra ID. Once there, in the left-hand pannel,
click on App Registrations. 
2. Then, click on New Registration in the top-middle of the screen.
   3. Give it a name
   4. Select Accounts only in this tenant (single tenant)
   5. Select Web for the platform, and give it a callback URL. In our case this will be ```http://localhost:8080/callback```
   6. Register the new App Registration
4. Next, we are going to update the App Registrations Manifest file. Click on Manifest located in the left-hand panel.
   5. change ```"accessTokenAcceptedVersion": null``` to ```"accessTokenAcceptedVersion": 2``` This is needed to ensure we are going to interact with the proper OIDC issuer. Save this change.
6. Next, create a secret. Click on Certificates and Secrets on the left-hand panel. Then click on Client Secrets. Create a new Secret, and store the value somewhere safe.
7. Now we need to grant API Permissions, The app registration will have one by default ```microsoft.graph.user.read```, we need to add one additional permission: ```microsoft.graph.profile``` this is needed so that we can configure a token claim to contain user UPN for logging and auditing purposes. Save these changes.
8. Lastly we need to create a new custom Token Configuration. Click on Token Configurations on the left-hand panel.
   9. Select Add optional claim
      10. Token Type, Select Access, and then scroll down until you see upn in the claims list. Select that, and then save.
