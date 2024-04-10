# Getting started with blazem

![demo image](./demo.png)

## Clone the repo

```
git clone https://github.com/JackElli/blazem.git
```
### We use `Docker` for our builds.
## To run blazem for development, run:
```
docker-compose up dev --build -d
```
- This will create a blazem instance, you can access it on any browser at address `localhost:5173`
- Once you are at the login screen, use username `JackTest` and password `test123` to log in.
- From there, you should be able to see your root folders.

(If an error occurs saying that the blazem_default network doesn't exist, just run the following:)

```
docker network create blazem-combined_default
```

### To add a document
First, add a folder, then click on the folder, then `Add data` and type in your value. (The key should automatically be generated for you)

###  To query
Go to the `Advanced search` tab and use 
```
SELECT all WHERE value LIKE 'a'
``` 
to select all documents that include the letter 'a'

## Functionality that isn't currently implemented

- Global search (the search bar at the top of the screen)
- Stats
- Import and export of data
- Multi-node deployment (was removed to focus on core product)
- Proper privated folders for multiple users
