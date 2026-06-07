import { useState } from "react";
import { Nav } from "./Nav";
import { Hero } from "./Hero";
import { Stats } from "./Stats";
import { PlatformTicker } from "./PlatformTicker";
import { Features } from "./Features";
import { HowItWorks } from "./HowItWorks";
import { Pricing } from "./Pricing";
import { FAQ } from "./FAQ";
import { CTASection } from "./CTASection";
import { Footer } from "./Footer";
import { AuthModal } from "./AuthModal";

/** Public marketing page shown to logged-out visitors. */
export function Landing() {
  const [authOpen, setAuthOpen] = useState(false);
  const openAuth = () => setAuthOpen(true);

  return (
    <div id="top">
      <Nav onAuth={openAuth} onKeys={openAuth} />
      <Hero onStart={openAuth} />
      <Stats />
      <PlatformTicker />
      <Features />
      <HowItWorks />
      <Pricing onChoose={openAuth} />
      <FAQ />
      <CTASection onStart={openAuth} />
      <Footer />
      <AuthModal open={authOpen} onClose={() => setAuthOpen(false)} />
    </div>
  );
}
