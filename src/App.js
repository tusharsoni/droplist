import React from "react";
import { H1 } from "baseui/typography";
import { Button } from "baseui/button";

function App() {
  return (
    <div>
      <H1>Hello, World!</H1>
      <Button onClick={() => alert("click")}>Hello</Button>
    </div>
  );
}

export default App;
