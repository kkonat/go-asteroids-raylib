#version 330

// Input vertex attributes (from vertex shader)
in vec2 fragTexCoord;
in vec4 fragColor;
in vec2 vertPos;

// Input uniform values
uniform sampler2D texture0;
uniform vec4 colDiffuse;
uniform vec2 pos, dir;
uniform float angle, size;


// Output fragment color
out vec4 finalColor;

// NOTE: Add here your custom variables

void main()
{
    // Texel color fetching from texture sampler
    vec4 texelColor = texture(texture0, fragTexCoord)*colDiffuse*fragColor;
    
    vec2 p = gl_FragCoord.xy - pos;
    float d = 1-length(p)/size;
    
    // Calculate final fragment color
    vec4 col = fragColor * d;
    finalColor = vec4(col.r, col.g, col.b, d*0.5);
}
