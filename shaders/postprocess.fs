#version 330

// Input vertex attributes (from vertex shader)
in vec2 fragTexCoord;
in vec4 fragColor;

// Input uniform values
uniform sampler2D texture0;
uniform float gamma;
uniform float iTime;
uniform vec2 iResolution;
uniform int glitch;

// Output fragment color
out vec4 finalColor;


// Tape noise
// https://www.shadertoy.com/view/MlfSWr

vec4 hash42(vec2 p){
    
	vec4 p4 = fract(vec4(p.xyxy) * vec4(443.8975,397.2973, 491.1871, 470.7827));
    p4 += dot(p4.wzxy, p4+19.19);
    return fract(vec4(p4.x * p4.y, p4.x*p4.z, p4.y*p4.w, p4.x*p4.w));
}


float hash( float n ){
    return fract(sin(n)*43758.5453123);
}

// 3d noise function (iq's)
float n( in vec3 x ){
    vec3 p = floor(x);
    vec3 f = fract(x);
    f = f*f*(3.0-2.0*f);
    float n = p.x + p.y*57.0 + 113.0*p.z;
    float res = mix(mix(mix( hash(n+  0.0), hash(n+  1.0),f.x),
                        mix( hash(n+ 57.0), hash(n+ 58.0),f.x),f.y),
                    mix(mix( hash(n+113.0), hash(n+114.0),f.x),
                        mix( hash(n+170.0), hash(n+171.0),f.x),f.y),f.z);
    return res;
}

//tape noise
float nn(vec2 p){


    float y = p.y;
    float s = iTime*2.;
    
    float v = (n( vec3(y*.01 +s, 			1., 1.0) ) + .0)
          	 *(n( vec3(y*.011+1000.0+s, 	1., 1.0) ) + .0) 
          	 *(n( vec3(y*.51+421.0+s, 	1., 1.0) ) + .0)   
        ;
    //v*= n( vec3( (fragCoord.xy + vec2(s,0.))*100.,1.0) );
   	v*= hash42(   vec2(p.x +iTime*0.01, p.y) ).x +.3 ;

    
    v = pow(v+.3, 1.);
	if(v<.7) v = 0.;  //threshold
    return v;
}

// ImageGlitcher
// https://www.shadertoy.com/view/MtXBDs
//2D (returns 0 - 1)
float random2d(vec2 n) { 
    return fract(sin(dot(n, vec2(12.9898, 4.1414))) * 43758.5453);
}

float randomRange (in vec2 seed, in float min, in float max) {
		return min + random2d(seed) * (max - min);
}

// return 1 if v inside 1d range
float insideRange(float v, float bottom, float top) {
   return step(bottom, v) - step(top, v);
}

//inputs
float AMT = 0.2; //0 - 1 glitch amount
float SPEED = 0.6; //0 - 1 speed
   
void main( )
{
    float time = floor(iTime /100.0 * SPEED * 60.0);    
        vec2 uv = fragTexCoord.xy;// / iResolution.xy;
        
        //copy orig
        vec3 outCol = texture(texture0, uv).rgb;
    if( glitch){
       
        //randomly offset slices horizontally
        float maxOffset = AMT/2.0;
        float count=1.0;
        for (float i = 0.0; i < 20.0 * AMT; i += 1.0) {
            float sliceY = random2d(vec2(time , 2345.0 + float(i)));
            float sliceH = random2d(vec2(time , 9035.0 + float(i))) * 0.25;
            float hOffset = randomRange(vec2(time , 9625.0 + float(i)), -maxOffset, maxOffset);
            vec2 uvOff = uv;
            uvOff.x += hOffset;
            if (insideRange(uv.y, sliceY, fract(sliceY+sliceH)) == 1.0 ){
                outCol += texture(texture0, uvOff).rgb;
                count += 1;
            }
        }
        outCol /= count;
        //do slight offset on one entire channel
        float maxColOffset = AMT/6.0;
        float rnd = random2d(vec2(time , 9545.0));
        vec2 colOffset = vec2(randomRange(vec2(time , 9545.0),-maxColOffset,maxColOffset), 
                        randomRange(vec2(time , 7205.0),-maxColOffset,maxColOffset));
        if (rnd < 0.33){
            outCol.r = texture(texture0, uv + colOffset).r;
            
        }else if (rnd < 0.66){
            outCol.g = texture(texture0, uv + colOffset).g;
            
        } else{
            outCol.b = texture(texture0, uv + colOffset).b;  
        }
        
        float linesN = 720.; //fields per seconds
        float one_y = iResolution.y / linesN; //field line
        uv = floor(uv*iResolution.xy/one_y)*one_y;
        outCol *=  (1.+4*nn(uv));
    }
    
 
    
	vec3 col = pow(outCol, vec3(1/gamma));
    
	finalColor = vec4(col,fragColor.a);
}

// void main()
// {
//     vec4 texelColor = texture(texture0, fragTexCoord)*fragColor;
//     vec3 col = texelColor.rgb;
//     //vec3 col = pow(texelColor.rgb, vec3(1/gamma));
//     //finalColor = vec4( col, texelColor.a);
//     finalColor = vec4( glitch( col, fragTexCoord), 1);
// }