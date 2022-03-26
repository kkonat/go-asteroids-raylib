#version 330

// Input vertex attributes (from vertex shader)
in vec2 fragTexCoord;
in vec4 fragColor;

// Input uniform values
uniform sampler2D texture0;

// Output fragment color
out vec4 finalColor;

// NOTE: Add here your custom variables

// NOTE: Render size values must be passed from code
const float renderWidth = 32;
const float renderHeight = 32;


    float Pi2 = 6.28318530718; // Pi*2
    
    // GAUSSIAN BLUR SETTINGS 
    float Directions = 16.0; // BLUR DIRECTIONS (Default 16.0 - More is better but slower)
    float Quality = 3.0; // BLUR QUALITY (Default 4.0 - More is better but slower)
    float Size = 1.1; // BLUR SIZE (Radius)
    // GAUSSIAN BLUR SETTINGS 
    
    vec2 Res = vec2(32.0,32.0);
    vec2 Radius = Size/Res;
 
void main()
{
   // Normalized pixel coordinates (from 0 to 1)
    vec2 uv = (fragTexCoord-vec2(16.0,16.0))/Res;
    // Pixel colour
    vec4 Color = texture(texture0, fragTexCoord);
    
    // Blur calculations
    for( float d=0.0; d<Pi2; d+=Pi2/Directions)
    {
		for(float i=1.0/Quality; i<=1.0; i+=1.0/Quality)
        {
			Color += texture( texture0, (uv+vec2(cos(d),sin(d)))*Res*i);		
        }
    }
    
    // Output to screen
    Color /= Quality * Directions -32;
    finalColor =  vec4(Color.r,Color.g,Color.b,1.0);
}