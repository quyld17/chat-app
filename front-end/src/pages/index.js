import Head from "next/head";
import SignIn from "./sign-in/index";

export default function Home() {
  return (
    <>
      <Head>
        <title>Chat App</title>
        <meta name="description" content="Real time chatting website" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <SignIn />
    </>
  );
}
