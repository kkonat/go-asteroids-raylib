#version 330

// Input vertex attributes (from vertex shader)
in vec2 fragTexCoord;
in vec4 fragColor;
in vec2 vertPos;

// Input uniform values
uniform vec2 iResolution;
uniform sampler2D texture0;
uniform vec4 colDiffuse;

// Output fragment color
out vec4 finalColor;

mat2 Rot(float a) {
    float s=sin(a),c=cos(a);
    return mat2(c,-s,s,c);
}
float Star(vec2 uv, float flare){
    float d = length(uv);
    float m = .203/d;
    
     float rays = max(0., 1.-abs(uv.x*uv.y*900.));        
     m += rays*0.2*flare;

     uv *= Rot(3.14159/4.);
     rays = max(0., 1.-abs(uv.x*uv.y*500.));
     m += rays*.14*flare;

    m *= smoothstep(1.,.2,d);
    return m;
}

void main()
{
    // Texel color fetching from texture sampler
    //vec4 texelColor = texture(texture0, fragTexCoord)*colDiffuse*fragColor;
    
    vec2 uv = fragTexCoord-0.5;
   uv *=2;
    vec3 col =vec3(0);
  //  uv *= 16;
    float l = length(uv);
    float d = 0.024/(l);
    float alpha = smoothstep(0.,1.,1-l);
    //col = fragColor.rgb* d;//smoothstep(0.0,1,d*d*d);
    col = fragColor.rgb * d*alpha;
    //col += fragColor.rgb *Star(uv,0.1);
    finalColor = vec4(col,1);
    //finalColor = fragColor*m; // Star(uv, 0.5);
}
