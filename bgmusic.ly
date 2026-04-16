\version "2.24.0"

\header {
  title = "Echoes at Midnight"
  subtitle = "Nocturne for Solo Piano"
  composer = "jason"
}

% --- GLOBAL MACROS ---
global = {
  \key a \minor
  \time 4/4
}

% --- RIGHT HAND ---
rightHand = \relative c'' {
  \clef treble
  \tempo "Lento con gran espressione" 4=56

  % === A SECTION: Lyrical and flowing ===
  \pp
  e4.( g8 ~ g a c b |
  a2 \acciaccatura { b16[ a] } g4. e8) |
  a4.( c8 ~ c b a g |
  f2. r4) |

  \mp
  e4.( g8 ~ g a c b |
  a2 \acciaccatura { b16[ a] } g4. e8) |
  d4.( f8 ~ f e d c |
  b2. \grace { c16[ b] } a8 b) |

  \mf
  c4.( e8 ~ e g f e |
  d4 f8 a ~ a2) |
  b,4.( d8 ~ d f e d |
  c1)\fermata \bar "||"

  % === B SECTION: Agitato (Driven and building) ===
  \tempo "Poco più mosso" 4=66
  \key c \major
  
  <<
    { \voiceOne 
      r4 e'4.\mf\< f8 g a |
      b2. c4\! |
      r4 a2.\ff ~ |
      a2. r4 |
    }
    \new Voice { \voiceTwo
      e,16\mf g c e g c e g  c, e g c e g e c |
      f, a c f a c f a  f, a c e c a f c |
      f,\ff a c f a c f a  f, a c e c a f c |
      e, g c e g c e g  e, g c e g e c g |
    }
  >>
  \oneVoice
  
  \key a \minor
  <<
    { \voiceOne 
      r4 e'4. f8 e d |
      c2. b4 |
    }
    \new Voice { \voiceTwo
      a,16\ff c e a c e a c  a, c e a c a e c |
      g, b d g b d g b  g, b d g b d b g |
    }
  >>
  \oneVoice
  
  % Descending cadential run
  \acciaccatura { b'8 } a4 \grace { g16[ f] } e8 d c b a g |
  f4\dim e8 d c b a g |
  f2.\> e8 d |
  e1\! \bar "||"

  % === A' SECTION: Ornamented Return ===
  \tempo "Tempo I" 4=56
  
  \p
  e'4.(\prall g8 ~ \turn g a c b |
  a2 \acciaccatura { b16[ a] } g4. e8) |
  a4.( \turn c8 ~ c b a g |
  f2. \acciaccatura { g16[ f] } e8 f) |

  \pp
  e4.( g8 ~ g a c b |
  a2 \acciaccatura { b16[ a] } g4. e8) |
  d4.( f8 ~ f e d c |
  
  % === LOOPING CADENCE ===
  % Ending on an E major chord creates a half-cadence, 
  % pulling the ear right back to the A minor start.
  b2. e4) \bar ":|."
}

% --- LEFT HAND ---
leftHand = \relative c, {
  \clef bass

  % A Section
  a8( e' a e' c e a, e' |
  a, e' a e' c e a, e') |
  f,( c' f c' a c f, c' |
  e,, b' e b' gis b e, b') |

  a,8( e' a e' c e a, e' |
  a, e' a e' c e a, e') |
  g,,8( d' g d' b d g, d' |
  c, g' c g' e g c, g') |

  c,,8( g' c g' e g c, g' |
  f,, c' f c' a c f, c') |
  g,,8( d' g d' b d g, d' |
  c, g' c g' e g c, g') |

  % B Section
  <c,, c'>8\mf[ <c c'> <c c'> <c c'>] <c c'>[ <c c'> <c c'> <c c'>] |
  <f, f'>[ <f f'> <f f'> <f f'>] <f f'>[ <f f'> <f f'> <f f'>] |
  <f, f'>\ff[ <f f'> <f f'> <f f'>] <f f'>[ <f f'> <f f'> <f f'>] |
  <c c'>[ <c c'> <c c'> <c c'>] <c c'>[ <c c'> <c c'> <c c'>] |

  <a, a'>\ff[ <a a'> <a a'> <a a'>] <a a'>[ <a a'> <a a'> <a a'>] |
  <e e'>[ <e e'> <e e'> <e e'>] <e e'>[ <e e'> <e e'> <e e'>] |
  <a, a'>4 r r2 |
  <d, d'>4 r r2 |
  <e, e'>2.\dim r4 |
  <a, a'>1 |

  % A' Section
  a8(\p e' a e' c e a, e' |
  a, e' a e' c e a, e') |
  f,( c' f c' a c f, c' |
  e,, b' e b' gis b e, b') |

  a,8(\pp e' a e' c e a, e' |
  a, e' a e' c e a, e') |
  g,,8( d' g d' b d g, d' |
  
  % === LOOPING CADENCE ===
  % E major arpeggio in the bass to match the melody
  e,8 b' e b' gis b e, e') \bar ":|."
}

% --- PEDALING ---
dynamics = {
  % A Section
  s1\sustainOn s2 s8 s\sustainOff s\sustainOn |
  s1 s4 s8 s\sustainOff s\sustainOn |
  s1 s2 s8 s\sustainOff s\sustainOn |
  s1 s4 s8 s\sustainOff s\sustainOn |
  
  s1 s2 s8 s\sustainOff s\sustainOn |
  s1 s4 s8 s\sustainOff s\sustainOn |
  s1 s2 s8 s\sustainOff s\sustainOn |
  s1\sustainOff |

  % B Section (Syncopated, driving pedal)
  s1\sustainOn s1\sustainOff\sustainOn |
  s1\sustainOff\sustainOn s1\sustainOff\sustainOn |
  s1\sustainOff\sustainOn s1\sustainOff\sustainOn |
  s1\sustainOff\sustainOn s1\sustainOff\sustainOn |
  
  s1\sustainOff s1\sustainOn |
  s4 s\sustainOff s2\sustainOn |
  s1\sustainOff s1\sustainOn |
  s1\sustainOff s1\sustainOn |
  s1\sustainOff s1 |

  % A' Section
  s1\sustainOn s2 s8 s\sustainOff s\sustainOn |
  s1 s4 s8 s\sustainOff s\sustainOn |
  s1 s2 s8 s\sustainOff s\sustainOn |
  s1 s4 s8 s\sustainOff s\sustainOn |
  
  s1 s2 s8 s\sustainOff s\sustainOn |
  s1 s4 s8 s\sustainOff s\sustainOn |
  s1 s2 s8 s\sustainOff s\sustainOn |
  
  % Hold the sustain pedal into the loop to ensure a seamless transition
  s1\sustainOn
}

% --- SCORE DEFINITION ---
\score {
  \new PianoStaff <<
    \repeat volta 2 {
      \new Staff = "RH" \rightHand
      \new Dynamics = "Dynamics" \dynamics
      \new Staff = "LH" \leftHand
    }
  >>
  \layout {
    \context {
      \PianoStaff
      \accepts Dynamics
    }
  }
  \midi {
    \context {
      \Voice
      \remove "Dynamic_performer"
    }
  }
}