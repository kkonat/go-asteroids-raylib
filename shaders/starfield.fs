#version 330
#define NUM_LAYERS 4.

// Input vertex attributes (from vertex shader)
in vec2 fragTexCoord;
in vec4 fragColor;

// Input uniform values
uniform sampler2D texture0;
uniform float time;
uniform vec2 iResolution;

// Output fragment color
out vec4 finalColor;

mat2 Rot(float a) {
    float s=sin(a),c=cos(a);
    return mat2(c,-s,s,c);
}
float Star(vec2 uv, float flare){
    float d = length(uv);
    float m = .103/d;
        
    float rays = max(0., 1.-abs(uv.x*uv.y*900.));
    m += rays*0.2*flare;

    uv *= Rot(3.14159/4.);
    rays = max(0., 1.-abs(uv.x*uv.y*1000.));
    m += rays*.14*flare;

    m *= smoothstep(1.,.2,d);
    return m;
}

float Hash21(vec2 p) {
    p = fract(p*vec2(123.34,456.21));
    p += dot(p,p+45.32);
    return fract(p.x*p.y);
}
vec3 StarLayer(vec2 uv) {
    vec3 col = vec3(0);
    vec2 gv = fract(uv)-.5;
    vec2 id = floor(uv);

//    if( gv.x>.48 || gv.y >.48) col.r=1.;
    for(int y=-1;y<=1;y++) {
            for(int x=-1;x<=1;x++) {
                vec2 offs = vec2(x,y);
                float n = Hash21(id+offs);
                float size = fract(n*345.32);
                
                float star = Star( gv-offs-vec2(n-0.5,fract(n*34.)-.5), 11.*smoothstep(.85,1.,size));
                vec3 color = sin(vec3(.2,.3,.9)*fract(n*2345.2)*6.238)*.5+.5;
                color = color*vec3(0.2+size,0+size,1.+size)*vec3(0.3,0.2,0.23);
                col += star*size*color;
            }
    }
    return col;
}
void main()
{
    float aspect = iResolution.x/iResolution.y;
    vec3 col = vec3(0);
    col = texture(texture0, fragTexCoord).rgb* vec3(0.5,0.3,0.8) - 0.1;
    vec2 uv = (fragTexCoord-0.5)*6.0;
    uv.y = uv.y / aspect;
    uv *= 5;
    float t = time*0.01;
    for( float i=0.; i<1.; i+=1./NUM_LAYERS) {
        float depth = fract(t*(i*0.6+1.));
        float offs = mix (20., .5, depth);
        float fade = depth*smoothstep(1.,.9,depth);
        col += StarLayer(uv-vec2(offs+i*454.3,0.))*fade;
    }
  
    finalColor = vec4(col,1.0)*0.7;
}