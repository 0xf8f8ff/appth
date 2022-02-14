import React, {useState } from 'react';  
import { AuthClient } from './proto/appth_grpc_web_pb';
import { UserRequest } from './proto/appth_pb';
import './App.css';

 // We create a client that connects to the api
 var client = new AuthClient("https://localhost:8080");

 function App() {
  // Create a const named status and a function called setStatus
  const [user, setUser] = useState(null);
  const requestUser = () => {
    var request = new UserRequest();
    request.setId(1);
    // use the client to send our request, the function that is passed
    // as the third param is a callback. 
    client.getUser(request, null, function(err, response) {
      if (err) {
        console.log("ERROR");
        console.log(err);
        return;
      }
      // serialize the response to an object 
      var user = response.toObject();
       setUser(user);
       console.log(response);
    }); 
  }

  return (
    <div className="App">
      <button onClick={
        () => {
          console.log('click');
          requestUser();
        }
      }>Fetch current user</button>
      <button onClick={
        () => {
          console.log('logging in...');
          window.location.assign("/auth/login");
        }
      }>Login</button>
      <button onClick={
        () => {
          console.log('redirecting...');
          window.location.assign("/auth/callback");
        }
      }>Redirect</button>
      <p>User: {user.getName()}</p>
    </div>
  );
}

export default App;
